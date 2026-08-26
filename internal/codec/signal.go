package codec

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

var ErrInvalidSignal = errors.New("invalid signal payload")

type SignalPacket struct {
	Version    byte
	Sequence   uint32
	SampleRate uint32
	Samples    []float64
	Checksum   uint32
}

func EncodeSignal(packet SignalPacket) ([]byte, error) {
	if packet.SampleRate == 0 || len(packet.Samples) == 0 {
		return nil, ErrInvalidSignal
	}
	if packet.Version == 0 {
		packet.Version = 1
	}
	buffer := make([]byte, 17+len(packet.Samples)*8)
	buffer[0] = packet.Version
	binary.BigEndian.PutUint32(buffer[1:5], packet.Sequence)
	binary.BigEndian.PutUint32(buffer[5:9], packet.SampleRate)
	binary.BigEndian.PutUint32(buffer[9:13], uint32(len(packet.Samples)))
	for i, sample := range packet.Samples {
		binary.BigEndian.PutUint64(buffer[13+i*8:21+i*8], math.Float64bits(sample))
	}
	checksum := checksumBytes(buffer[:13+len(packet.Samples)*8])
	binary.BigEndian.PutUint32(buffer[13+len(packet.Samples)*8:], checksum)
	return buffer, nil
}

func DecodeSignal(data []byte) (SignalPacket, error) {
	if len(data) < 17 {
		return SignalPacket{}, ErrInvalidSignal
	}
	count := int(binary.BigEndian.Uint32(data[9:13]))
	if count <= 0 {
		return SignalPacket{}, ErrInvalidSignal
	}
	// 校验声明的样本数量没有超出实际可用载荷，防止损坏的长度字段导致越界读取。
	if len(data) < 13+count*8+4 {
		return SignalPacket{}, fmt.Errorf("%w: sample count %d exceeds payload of %d bytes", ErrInvalidSignal, count, len(data))
	}
	expected := binary.BigEndian.Uint32(data[len(data)-4:])
	actual := checksumBytes(data[:len(data)-4])
	if expected != actual {
		return SignalPacket{}, errors.New("signal checksum mismatch")
	}
	packet := SignalPacket{Version: data[0], Sequence: binary.BigEndian.Uint32(data[1:5]), SampleRate: binary.BigEndian.Uint32(data[5:9]), Samples: make([]float64, count), Checksum: expected}
	for i := range packet.Samples {
		packet.Samples[i] = math.Float64frombits(binary.BigEndian.Uint64(data[13+i*8 : 21+i*8]))
	}
	return packet, nil
}

func checksumBytes(data []byte) uint32 {
	var checksum uint32 = 2166136261
	for _, value := range data {
		checksum ^= uint32(value)
		checksum *= 16777619
	}
	return checksum
}

func EncodeSamplesCSV(samples []float64) string {
	parts := make([]string, len(samples))
	for i, sample := range samples {
		parts[i] = strconv.FormatFloat(sample, 'f', 6, 64)
	}
	return strings.Join(parts, ",")
}

func DecodeSamplesCSV(raw string) ([]float64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, ErrInvalidSignal
	}
	parts := strings.Split(raw, ",")
	out := make([]float64, len(parts))
	for i, part := range parts {
		value, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, err
		}
		out[i] = value
	}
	return out, nil
}

func ReadSignalLines(reader io.Reader) ([]SignalPacket, error) {
	scanner := bufio.NewScanner(reader)
	packets := make([]SignalPacket, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		samples, err := DecodeSamplesCSV(line)
		if err != nil {
			return packets, err
		}
		packets = append(packets, SignalPacket{Version: 1, SampleRate: 1000, Sequence: uint32(len(packets)), Samples: samples})
	}
	return packets, scanner.Err()
}

func WriteSignalLines(writer io.Writer, packets []SignalPacket) error {
	for _, packet := range packets {
		if _, err := io.WriteString(writer, EncodeSamplesCSV(packet.Samples)+"\n"); err != nil {
			return err
		}
	}
	return nil
}
