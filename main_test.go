package main

import (
  "testing"
  "os"
)

func TestSanitizeSpaces (t *testing.T) {
  loadWords()
  actual := sanitizeSpaces("sh   it")
  expected := "****"
  if actual != expected {
    t.Fatalf("TestSanitizeSpaces expected %s got %s", expected, actual)
  }
}

func TestSanitizeSymbols (t *testing.T) {
  loadWords()
  actual := sanitizeSymbols("$hit")
  expected := "****"
  if actual != expected {
    t.Fatalf("TestSanitizeSymbols expected %s got %s", expected, actual)
  }
}

func TestSanitizeConcurrentLetters (t *testing.T) {
  loadWords()
  actual := sanitizeConcurrentLetters("shiiiiiit")
  expected := "****"
  if actual != expected {
    t.Fatalf("TestSanitizeConcurrentLetters expected %s got %s", expected, actual)
  }
}

func TestSanitize (t *testing.T) {
  loadWords()
  actual := sanitize("$hiiiiiit fuuuck man")
  expected := "**** **** man"
  if actual != expected {
    t.Fatalf("TestSanitize expected %s got %s", expected, actual)
  }
}

func TestMain (m *testing.M) {
  loadWords()
  os.Exit(m.Run())
}
