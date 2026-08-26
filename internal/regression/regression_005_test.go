package regression

import (
    "context"
    "database/sql"
    "errors"
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

type failingArchive struct{ err error }
func (f failingArchive) Write(context.Context, model.Survey, []model.EchoProfile) (model.Archive, error) { return model.Archive{}, f.err }

func TestBug05_ArchiveWriterErrorChain(t *testing.T) {
    _, repo := fixture(t)
    diskErr := errors.New("disk quota")
    svc := service.New(repo, service.Config{Archive: failingArchive{err: diskErr}})
    defer svc.Close()
    if err := svc.CloseSurvey(context.Background(), "survey-1"); err != nil { t.Fatal(err) }
    err := svc.ArchiveSurvey(context.Background(), "survey-1")
    if !errors.Is(err, service.ErrArchiveFailed) || !errors.Is(err, diskErr) { t.Fatalf("error identity lost: %v", err) }
}

func TestBug05_MissingArchiveDistinctFromFailure(t *testing.T) {
    svc, _ := fixture(t)
    _, err := svc.GetArchive(context.Background(), "survey-missing")
    if errors.Is(err, service.ErrArchiveFailed) { t.Fatal("missing archive should not look like writer failure") }
    if !errors.Is(err, service.ErrArchiveNotFound) { t.Fatalf("want ErrArchiveNotFound, got %v", err) }
}

func TestBug05_StoreErrorChainPreserved(t *testing.T) {
    _, repo := fixture(t)
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    _, err := repo.GetArchive(ctx, "survey-1")
    if !errors.Is(err, context.Canceled) { t.Fatalf("store error identity lost: %v", err) }
}