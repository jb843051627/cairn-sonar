package regression

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "path/filepath"
    "sync"
    "strings"
    "testing"
    "time"

    "cairn-sonar/internal/api"
    "cairn-sonar/internal/clock"
    "cairn-sonar/internal/codec"
    "cairn-sonar/internal/model"
    "cairn-sonar/internal/rules"
    "cairn-sonar/internal/service"
    "cairn-sonar/internal/store"
    "cairn-sonar/internal/worker"
)

var _ = errors.Is
var _ = api.ErrBodyTooLarge
var _ = codec.ErrInvalidSignal
var _ = rules.ErrInvalidWindow
var _ = sync.Once{}
var _ = strings.Builder{}
var _ = worker.ErrBatchCancelled

func fixture(t *testing.T) (*service.Service, *store.Repository) {
    t.Helper()
    repo, err := store.Open(filepath.Join(t.TempDir(), "survey.db"))
    if err != nil { t.Fatal(err) }
    now := clock.NewFake(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
    svc := service.New(repo, service.Config{Clock: now})
    chamber := model.Chamber{ID: "chamber-1", Name: "North wall", SiteCode: "CAIRN-A", DepthM: 18, CreatedAt: now.Now()}
    if err := repo.PutChamber(context.Background(), chamber); err != nil { t.Fatal(err) }
    instrument := model.Instrument{ID: "instrument-1", Serial: "SN-1", Firmware: "1.4", FrequencyHz: 40000, Enabled: true, LastCalibrated: now.Now()}
    if err := repo.PutInstrument(context.Background(), instrument); err != nil { t.Fatal(err) }
    survey := model.Survey{ID: "survey-1", ChamberID: chamber.ID, Lead: "operator-1", Status: model.SurveyPlanned, StartedAt: now.Now()}
    if err := svc.CreateSurvey(context.Background(), survey); err != nil { t.Fatal(err) }
    if err := svc.StartSurvey(context.Background(), survey.ID); err != nil { t.Fatal(err) }
    t.Cleanup(func() { svc.Close(); _ = repo.Close() })
    return svc, repo
}

func pulse(surveyID string) model.Pulse {
    return model.Pulse{ID: "pulse-1", SurveyID: surveyID, InstrumentID: "instrument-1", Sequence: 1, EmittedAt: time.Date(2026, 1, 2, 3, 5, 0, 0, time.UTC), DurationMS: 20, Samples: []float64{-12, -10, 72, 34, 12}}
}

func failTx(tx *sql.Tx) error {
    if _, err := tx.Exec("INSERT INTO events(survey_id,kind,payload,created_at) VALUES(?,?,?,?)", "survey-1", "probe", "partial", "2026-01-02T03:04:05Z"); err != nil { return err }
    return errors.New("stop after partial event")
}

func anomaly() model.Anomaly {
    return model.Anomaly{ID: "anomaly-1", SurveyID: "survey-1", EchoID: "pulse-1", Kind: "delamination", Severity: 2, Evidence: []float64{-18, -12, 66, 31}}
}

func TestBug08_EvidenceArrayIsolation(t *testing.T) {
    svc, _ := fixture(t)
    if err := svc.RaiseAnomaly(context.Background(), anomaly()); err != nil { t.Fatal(err) }
    first, err := svc.ListAnomalies(context.Background(), "survey-1"); if err != nil { t.Fatal(err) }
    first[0].Evidence[0] = 999
    second, err := svc.ListAnomalies(context.Background(), "survey-1"); if err != nil { t.Fatal(err) }
    if second[0].Evidence[0] == 999 { t.Fatal("cached evidence was mutated through caller slice") }
}

func TestBug08_RaisedInputDoesNotPolluteCache(t *testing.T) {
    svc, _ := fixture(t)
    a := anomaly()
    if err := svc.RaiseAnomaly(context.Background(), a); err != nil { t.Fatal(err) }
    a.Evidence[0] = 999
    items, err := svc.ListAnomalies(context.Background(), "survey-1"); if err != nil { t.Fatal(err) }
    if items[0].Evidence[0] == 999 { t.Fatal("raised anomaly input leaked into cached evidence") }
}

func TestBug08_ReviewReflectsInCachedList(t *testing.T) {
    svc, _ := fixture(t)
    if err := svc.RaiseAnomaly(context.Background(), anomaly()); err != nil { t.Fatal(err) }
    if err := svc.ReviewAnomaly(context.Background(), "anomaly-1", model.AnomalyAccepted, "ranger-1", "confirmed"); err != nil { t.Fatal(err) }
    items, err := svc.ListAnomalies(context.Background(), "survey-1"); if err != nil { t.Fatal(err) }
    if items[0].State != model.AnomalyAccepted { t.Fatalf("reviewed state not reflected in list: %s", items[0].State) }
}

func TestBug08_ConcurrentReadsAreRaceFree(t *testing.T) {
    svc, _ := fixture(t)
    if err := svc.RaiseAnomaly(context.Background(), anomaly()); err != nil { t.Fatal(err) }
    stop := make(chan struct{})
    var wg sync.WaitGroup
    for i := 0; i < 4; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-stop:
                    return
                default:
                    if _, err := svc.ListAnomalies(context.Background(), "survey-1"); err != nil { t.Error(err); return }
                }
            }
        }()
    }
    for i := 0; i < 3; i++ {
        a := anomaly()
        a.ID = fmt.Sprintf("anomaly-%d", i+2)
        if err := svc.RaiseAnomaly(context.Background(), a); err != nil { t.Fatal(err) }
    }
    close(stop)
    wg.Wait()
}