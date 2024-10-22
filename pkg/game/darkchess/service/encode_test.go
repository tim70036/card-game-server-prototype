package service

import (
	"fmt"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	gameMetaUid := "01140a00-0000-0000-0500-013880000000"
	rawGameMetaUid := strings.ReplaceAll(gameMetaUid, "-", "")
	if (len(rawGameMetaUid)) < 18 {
		t.Fatal(fmt.Errorf("invalid rawGameMetaUid %s, too short", rawGameMetaUid))
	}

	totalRound, err := strconv.ParseInt(rawGameMetaUid[0:2], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse totalRound from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("totalRound: %d", totalRound)

	turnSecond, err := strconv.ParseInt(rawGameMetaUid[2:4], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse turnSecond from string %s: %w", rawGameMetaUid[4:6], err))
	}
	t.Logf("turnSecond: %d->%v", turnSecond, time.Duration(turnSecond)*time.Second)

	extraTurnSecond, err := strconv.ParseInt(rawGameMetaUid[4:6], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse extraTurnSecond from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("extraTurnSecond: %d->%v", extraTurnSecond, time.Duration(extraTurnSecond)*time.Second)

	isCaptureRevealPieces, err := strconv.ParseInt(rawGameMetaUid[8:9], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse isCaptureRevealPieces from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("isCaptureRevealPieces: %d->%v", isCaptureRevealPieces, isCaptureRevealPieces == 1)

	isCaptureUnrevealPiece, err := strconv.ParseInt(rawGameMetaUid[9:10], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse isCaptureUnrevealPiece from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("isCaptureUnrevealPiece: %d->%v", isCaptureUnrevealPiece, isCaptureUnrevealPiece == 1)

	isCaptureUnrevealPieces, err := strconv.ParseInt(rawGameMetaUid[10:11], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse isCaptureUnrevealPieces from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("isCaptureUnrevealPieces: %d->%v", isCaptureUnrevealPieces, isCaptureUnrevealPieces == 1)

	hasRookRules, err := strconv.ParseInt(rawGameMetaUid[11:12], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse hasRookRules from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("hasRookRules: %d->%v", hasRookRules, hasRookRules == 1)

	hasBishopRules, err := strconv.ParseInt(rawGameMetaUid[12:13], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse hasBishopRules from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("hasBishopRules: %d->%v", hasBishopRules, hasBishopRules == 1)

	waterPct, err := strconv.ParseInt(rawGameMetaUid[16:18], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse waterPct from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("waterPct: %d", waterPct)

	anteAmount, err := strconv.ParseInt(rawGameMetaUid[20:25], 16, 0)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to parse anteAmount from string %s: %w", rawGameMetaUid[4:6], err))
	}

	t.Logf("anteAmount: %d", anteAmount)
	t.Logf("EnterLimit: %d", int(anteAmount)*constant.BuddyEnterLimitWeight)
}
