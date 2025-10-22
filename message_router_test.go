package main

import (
	"reflect"
	"testing"
)

func handlerPointer(fn RouterHandler) uintptr {
	if fn == nil {
		return 0
	}
	return reflect.ValueOf(fn).Pointer()
}

func TestMessageRouterDetermine(t *testing.T) {
	t.Parallel()

	router := NewMessageRouter()

	tests := []struct {
		name        string
		input       string
		wantHandler RouterHandler
		wantErr     bool
	}{
		{
			name:        "hello command",
			input:       "hello",
			wantHandler: helloHandler,
		},
		{
			name:        "hello command with extra text",
			input:       "hello there",
			wantHandler: helloHandler,
		},
		{
			name:        "hello command with leading mention",
			input:       "<@U123> hello",
			wantHandler: helloHandler,
		},
		{
			name:        "help command with multiple mentions",
			input:       "<@U123> <@U456> help me",
			wantHandler: helpHandler,
		},
		{
			name:    "unknown command",
			input:   "unknown",
			wantErr: true,
		},
		{
			name:    "empty message",
			input:   "",
			wantErr: true,
		},
		{
			name:    "blank spaces",
			input:   "   ",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler, err := router.Determine(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got none", tc.input)
				}
				if handler != nil {
					t.Fatalf("expected nil handler for input %q, got non-nil", tc.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if handler == nil {
				t.Fatalf("expected handler for input %q, got nil", tc.input)
			}

			gotPointer := handlerPointer(handler)
			wantPointer := handlerPointer(tc.wantHandler)
			if gotPointer != wantPointer {
				t.Fatalf("unexpected handler for input %q: got %v, want %v", tc.input, gotPointer, wantPointer)
			}
		})
	}
}
