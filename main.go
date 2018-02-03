package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "log"
    "math/rand"
    "os"
    "strings"
    "time"
)

type result struct {
    CorrectAnswers    int
    NumberOfQuestions int
}

type question struct {
    Question string
    Answer   string
}

func main() {
    csvfile := flag.String("csvfile", "problems.csv", "csv file with question")
    timeOut := flag.Int("timeout", 5, "timeout for answers")
    sort := flag.Bool("sortrandom", false, "sort randomly")
    flag.Parse()
    questions, err := os.Open(*csvfile)
    if err != nil {
        log.Fatal(err)
    }
    defer questions.Close()

    reader := csv.NewReader(questions)
    lines, err := reader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    qs := parseLines(lines)

    if *sort {
        randomSort(qs)
    }

    result := quiz(qs, *timeOut)
    fmt.Printf("You have %d of correct answers for total of %d question \n",
        result.CorrectAnswers, result.NumberOfQuestions)
}
func randomSort(qs []question) {
    for i := range qs {
        j := rand.Intn(i + 1)
        qs[i], qs[j] = qs[j], qs[i]
    }
}
func parseLines(lines [][]string) []question {
    ret := make([]question, len(lines))
    for i, line := range lines {
        ret[i] = question{
            Question: strings.TrimSpace(line[0]),
            Answer:   strings.TrimSpace(line[1]),
        }
    }
    return ret
}
func quiz(qs []question, timeout int) result {
    correct := 0
    timer := time.NewTimer(time.Duration(timeout) * time.Second)
    for _, p := range qs {
        select {
        case <-timer.C:
            return result{NumberOfQuestions: len(qs), CorrectAnswers: correct}
        default:
            fmt.Print(p.Question + " :")
            text := ""
            fmt.Scanf("%s\n", &text)
            if strings.TrimSpace(text) == p.Answer {
                correct += 1
            }
        }
    }
    return result{NumberOfQuestions: len(qs), CorrectAnswers: correct}
}
