package main
import (
    "fmt"
    "math"
)
func printMaxTag(posUnigrams []PosUniGram,previous []float64) string {
    maxVal := math.Log(0)
    maxTag := posUnigrams[0].PosTag
    for idx,j := range posUnigrams {
        if(previous[idx] != 0 && math.Exp(previous[idx]) > maxVal) {
            maxVal = math.Exp(previous[idx])
            maxTag = j.PosTag
        }
    }
    return maxTag
}

func processWord(word string,posUnigrams []PosUniGram, previous []float64) []float64 {
    next := make([]float64,82)
    possibleTags := getWordPosgramCount(word)
    for idx,j := range posUnigrams {
        for _,k := range possibleTags {
            if (k.PosTag == j.PosTag) {
                //next[idx] = math.Log(float64(k.Count)) - math.Log(float64(j.Count)) + previous[idx]
                next[idx] = float64(k.Count)*previous[idx]/float64(j.Count)
            }
        }
    }
    //fmt.Printf("%v",next)
    return next
}

func processTag(tag string,posUnigrams []PosUniGram,val float64,previous []float64) []float64 {
    next := make([]float64,82)
    i := 0
    for idx,j := range posUnigrams {
            i += 1
            count := getPosgram( tag+" "+j.PosTag)
            //fmt.Printf("%.2f %.2f\n",math.Log(float64(count)) - math.Log(float64(j.Count)) + val,previous[idx])
            //next[idx] = math.Max(math.Log(float64(count)) - math.Log(float64(j.Count)) + val,previous[idx])
            next[idx] = math.Max(float64(count)*val/float64(j.Count),previous[idx])
    }
    //fmt.Printf("%v\n",i)
    return next
}

func fastProcessTag(tag string,posUnigrams []PosUniGram,val float64,previous []float64) []float64 {
    next := make([]float64,82)
    list := make([]string,82)
    for idx,j := range posUnigrams {
        list[idx] =  tag+" "+j.PosTag
    }
    result := getAllPosgram(list)
    i := 0
    for idx,j := range posUnigrams {
        change := true
        for _,k := range result {
            if (k.Ngram == (tag+" "+j.PosTag)) {
                i += 1
                next[idx] = math.Max(float64(k.Count)*val/float64(j.Count),previous[idx])
                change = false
            }
        }
        if(change) {
            next[idx] = math.Max(0,previous[idx])
        }
    }
    //fmt.Printf("%v\n",next)
    return next
}

func viterbi(sentence string) {
    previous := make([]float64,82)
    next := make([]float64,82)
    regularexp := GetRegex()
    tokens := ProcessLine(sentence,regularexp)
    tokens = append(tokens,"<\\s>")
    values := getAllPosUnigrams()
    tag := "starts"
    for idx,j := range values {
        count := getPosgram(tag + " " + j.PosTag)
        //previous[idx] = math.Log(float64(count)) - math.Log(float64(j.Count))
        previous[idx] = float64(count)/float64(j.Count)
    }
    previous = processWord(tokens[1],values,previous)
    fmt.Printf("POS:%s\n",printMaxTag(values,previous))
    for _,i := range tokens[2:len(tokens)-1] {
        for idx,j := range values {
            next = fastProcessTag(j.PosTag,values,previous[idx],next)
            //next = processTag(j.PosTag,values,previous[idx],next)
            //fmt.Printf("%v\n",next)
            //fmt.Printf("============================\n")
        }
        next = processWord(i,values,next)
        fmt.Printf("POS:%s\n",printMaxTag(values,next))
        //fmt.Printf("%v",next)
        previous = next
        next = make([]float64,82)
    }
    // for idx,_ := range previous {
    //     fmt.Printf("%s %.2f\n",values[idx].PosTag,previous[idx])
    // }
}
