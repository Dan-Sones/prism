package prismhash

import (
	"testing"
)

func TestGetNumberLinePositionDeterminism(t *testing.T) {
	a := GetNumberLinePositionForUserAndExperiment("user_123", "exp_a", "salt_1")
	b := GetNumberLinePositionForUserAndExperiment("user_123", "exp_a", "salt_1")

	if a != b {
		t.Errorf("same inputs should deliver same results, got %d and %d", a, b)
	}
}

func TestGetNumberLinePositionInRange(t *testing.T) {
	users := []string{"user_1", "user_2", "user_3", "user_4", "user_5", "alice", "bob", ""}

	for _, u := range users {
		pos := GetNumberLinePositionForUserAndExperiment(u, "exp_a", "salt_1")
		if pos < 0 || pos > 99 {
			t.Errorf("expected between 0 and 99 got %d", pos)
		}
	}
}

func TestGetNumberLinePositionUserIdAffectsOutput(t *testing.T) {
	a := GetNumberLinePositionForUserAndExperiment("user_123", "exp_a", "salt_1")
	b := GetNumberLinePositionForUserAndExperiment("user_456", "exp_a", "salt_1")

	if a == b {
		t.Error("Different user ids should produce different positions")
	}
}

func TestGetNumberLinePositionExperimentKeyAffectsOutput(t *testing.T) {
	a := GetNumberLinePositionForUserAndExperiment("user_123", "exp_a", "salt_1")
	b := GetNumberLinePositionForUserAndExperiment("user_123", "exp_b", "salt_1")

	if a == b {
		t.Error("Different exp keys should produce different positions")
	}
}

func TestGetNumberLinePositionSaltAffectsOutput(t *testing.T) {
	a := GetNumberLinePositionForUserAndExperiment("user_123", "exp_a", "salt_1")
	b := GetNumberLinePositionForUserAndExperiment("user_123", "exp_a", "salt_2")

	if a == b {
		t.Error("Different salts should produce different positions")
	}
}

func TestGetNumberLinePositionStability(t *testing.T) {
	cases := []struct {
		userId        string
		experimentKey string
		uniqueSalt    string
		expected      int32
	}{
		{"user_123", "exp_a", "salt_1", 54},
		{"user_456", "exp_a", "salt_1", 73},
		{"user_789", "checkout_redesign", "abc123", 52},
	}

	for _, c := range cases {
		actual := GetNumberLinePositionForUserAndExperiment(c.userId, c.experimentKey, c.uniqueSalt)
		if actual != c.expected {
			t.Errorf("expected %d, got %d", c.expected, actual)
		}
	}
}
