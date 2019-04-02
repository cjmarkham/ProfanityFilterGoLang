package main

import  (
  "fmt"
  "os"
  "io/ioutil"
  "encoding/json"
  "github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
  "strings"
)

type symbolMapping struct {
  Symbol string
  Letter rune
}

type sanitizeResponse struct {
  IsSanitized bool
  Sanitized string
  Original string
}

var (
  profaneWords []string
  symbolMappings = []symbolMapping {
    symbolMapping {
      "$", 's',
    },
    symbolMapping {
      "!", 'i',
    },
  }
  sanitized string
)

func main () {
  loadWords()

  // os.Args first arg is the program name
  if len(os.Args) == 1 {
    panic("Please provide a word to sanitize")
  }

  sanitized := sanitize(os.Args[1])
  fmt.Println(sanitized)
}

func sanitize (s string) string {
  sanitized = sanitizeSpaces(s)
  sanitized = sanitizeConcurrentLetters(sanitized)
  sanitized = sanitizeSymbols(sanitized)

  return sanitized
}

func loadWords () {
  path := "./words.json"
  jsonFile, err := os.Open(path)

  if (err != nil) {
    fmt.Printf(err.Error())
    return
  }

  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)
  json.Unmarshal([]byte(byteValue), &profaneWords)
}

func sanitizeSpaces (s string) string {
  r := pcre.MustCompile("(\\s)\\1+", 0)
  re := r.MatcherString(s, 0)

  stringWithoutConcurrency := s

  if re.Matches() {
    match := re.GroupString(1)
    stringWithoutConcurrency = string(r.ReplaceAll([]byte(s), []byte(match), 0))
  }

  // Loop through all words in our profanity dictionary
  for i := 0; i < len(profaneWords); i++ {
    word := profaneWords[i]
    // Get all varients of this word with spaces
    spacedVariants := addSpacesToWord(word)

    for j := 0; j < len(spacedVariants); j++ {
      spaced := spacedVariants[j]

      letters := strings.Split(string(stringWithoutConcurrency), "")

      for index := 0; index < len(letters); index++ {
        letter := letters[index]
        // We need to replace any symbols in the string
        // This is pretty much the same as the symbol sanitization below
        // but we need to do this here too as that method is called after this
        // and relies on space removal
        for l := range symbolMappings {
          // If there is a replacement, replace the letter in the word
          // with the replacement
          // We dont want to directly replace the word as this symbol
          // could be a legit one
          if symbolMappings[l].Symbol == letter {
            // If the symbol is at the end of the word, leave it
            // This is an assumption that the last symbol is a legit one (!)
            if index == len(stringWithoutConcurrency) - 1 {
              break
            }

            // Replace the symbol with the letter in the string at this index
            runes := []rune(stringWithoutConcurrency)
            runes[index] = symbolMappings[l].Letter
            stringWithoutConcurrency = string(runes)
          }
        }
      }

      re := pcre.MustCompile(spaced, 0)
      reg := re.MatcherString(stringWithoutConcurrency, 0)

      if reg.Matches() {
        matched := strings.Replace(reg.GroupString(0), " ", "", -1)
        stringWithoutConcurrency = string(
          re.ReplaceAll(
            []byte(stringWithoutConcurrency), []byte(strings.Repeat("*", len(matched))), 0,
          ),
        )
      }
    }
  }

  return stringWithoutConcurrency
}

func sanitizeSymbols (s string) string {
  words := strings.Split(s, " ")

  for i := 0; i < len(words); i++ {
    letters := strings.Split(s, "")
    newWord := ""

    for j := 0; j < len(letters); j++ {
      for k := range symbolMappings {
        // If there is a replacement, replace the letter in the word
        // with the replacement
        // We dont want to directly replace the word as this symbol
        // could be a legit one
        if symbolMappings[k].Symbol == letters[j] {
          // If the symbol is at the end of the word, leave it
          if j == len(s) - 1 {
            break
          }

          // If we have already replaced a letter in this word, we
          // need to use the word stored in new_word
          // We need to replace just the symbol at this index as there
          // may be legit symbols in this word that are the same
          if newWord == "" {
            newWord = s
          }

          // Replace the symbol with the letter in the string at this index
          runes := []rune(newWord)
          runes[k] = symbolMappings[k].Letter
          newWord = string(runes)
        }
      }
    }

    // No letters were replaced
    if newWord == "" {
      break
    }

    sanitizeCheck := sanitizeWord(strings.ToLower(newWord))

    if sanitizeCheck.IsSanitized {
      s = strings.Replace(s, s, sanitizeCheck.Sanitized, -1)
    }
  }

  return s
}

func sanitizeConcurrentLetters (s string) string {
  words := strings.Split(s, " ")

  for i := 0; i < len(words); i++ {
    r := pcre.MustCompile("(\\w)\\1+", 0)
    re := r.MatcherString(words[i], 0)

    wordWithoutConcurrency := words[i]

    if re.Matches() {
      matched := re.GroupString(1)
      wordWithoutConcurrency = string(r.ReplaceAll([]byte(words[i]), []byte(matched), 0))
    }

    sanitizeCheck := sanitizeWord(wordWithoutConcurrency)

    if sanitizeCheck.IsSanitized {
      s = strings.Replace(s, words[i], sanitizeCheck.Sanitized, -1)
    }
  }

  return s
}

func sanitizeWord (w string) sanitizeResponse {
  isSanitized := false

  // Remove concurrent letters
  // If multiple symbols were used they will be replaced
  // with letters and there will be concurrency
  r := pcre.MustCompile("(\\w)\\1+", 0)
  re := r.MatcherString(w, 0)

  wordWithoutConcurrency := w

  if re.Matches() {
    matched := re.GroupString(1)
    wordWithoutConcurrency = string(r.ReplaceAll([]byte(w), []byte(matched), 0))
  }

  profanityMatch := ""
  for i := 0; i < len(profaneWords); i++ {
    if profaneWords[i] == wordWithoutConcurrency {
      profanityMatch = profaneWords[i]
      break
    }
  }

  sanitized := ""

  if profanityMatch != "" {
    isSanitized = true
    re := pcre.MustCompile(profanityMatch, 0)

    sanitized = string(
      re.ReplaceAll(
        []byte(w), []byte(strings.Repeat("*", len(profanityMatch))), 0,
      ),
    )
  }

  return sanitizeResponse {
    isSanitized,
    sanitized,
    wordWithoutConcurrency,
  }
}

func addSpacesToWord (word string) []string {
  if len(word) == 1 {
    return []string { word }
  }

  spacedWords := []string { }
  firstChar := string(word[0])
  recurse := addSpacesToWord(word[1:len(word)])

  for i := 0; i < len(recurse); i++ {
    spacedWords = append(spacedWords, firstChar + recurse[i], firstChar + ` ` + recurse[i])
  }

  return spacedWords
}

