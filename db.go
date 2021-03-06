package main

import (
	"fmt"
    "strings"
    "time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Tokens struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Tokens      []string
	//count     int
	Timestamp time.Time
}

type Ngrams struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Ngram      string
	Count     int
	Timestamp time.Time
}

type PosWordGram struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	PosTag      string		"posTag"
	Word		string
	Count     int
}

type ChunkPosGram struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	ChunkTag      string		"chunk_tag"
	WordPos		string			"word_pos"
	Count     int
}

type PosUniGram struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	PosTag      string		"posTag"
	Count     int
}

func InsertTokens(tokens []string) {
    maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    c := session.DB("nlprokz").C("brownCorpus")
    err = c.Insert(&Tokens{Tokens: tokens, Timestamp: time.Now()})

	if err != nil {
		panic(err)
	}
}

func InsertWordPosgram(tokens []string) {
    maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C("wordPosgram")
    //pos := session.DB("nlprokz").C("posUnigram")
	for j := 0; j < len(tokens); j++ {

        var temp = strings.Split(tokens[j],"/")
        if(len(temp) > 1) {
            posTag := temp[len(temp)-1]
            if(strings.Index(posTag,"+") != -1 && strings.Index(posTag,"+") != (len(posTag)-1)  ) {
                temp = strings.Split(posTag,"+")
                if(len(temp[1]) > 1 || temp[1] != "*") {
                    pos.Upsert(bson.M{"posTag":temp[1]},bson.M{"$inc": bson.M{"count": 1}})
                }
                posTag = temp[0]
            }
            if(strings.Contains(posTag,"-")) {
                posTag = strings.Split(posTag,"-")[0]
            }
            word := temp[0]
            if(len(temp) > 2) {
                word = strings.Join(temp[0:len(temp)-1],"/")
                fmt.Println("%s------------%s\n",tokens,word)
            }
            // if(strings.Contains(posTag,"-") ||strings.Contains(posTag,"+") || len(posTag) < 1 || posTag == "*") {
            //     fmt.Printf("%#v\n",tokens)
            // }

            pos.Upsert(bson.M{"posTag":posTag,"word":word},bson.M{"$inc": bson.M{"count": 1}})
        }
	}
}

func InsertChunkgram(chunk_list [][]string, posgrams []string) {
	maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    chunkgrams := session.DB("nlprokz").C("wordChunkgrams")
	chunks := session.DB("nlprokz").C("wordChunks")
	// fmt.Printf("%d\n",len(posgrams))
	fmt.Printf("%d\n",len(chunk_list))
	for i,_ := range chunk_list[0:len(chunk_list)-1] {
		// fmt.Printf(posgrams[i])
		// fmt.Printf(chunk_list[i][0])
		// fmt.Printf("==============\n")
		chunkgrams.Upsert(bson.M{"ngram":chunk_list[i+1][2]+" "+chunk_list[i][0]+" "+chunk_list[i+1][0]},bson.M{"$inc": bson.M{"count": 1}})
		chunks.Upsert(bson.M{"ngram":chunk_list[i+1][2]+" "+chunk_list[i][0]},bson.M{"$inc": bson.M{"count": 1}})
	}
}

func InsertChunkNgram(tokens []string, n int) {
    maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    chunk := session.DB("nlprokz").C("chunkngram")
	for j := 0; j < (len(tokens) - n); j++ {
		chunkSeq := make([]string,0)
		for a := 0; a < n; a++ {
			chunkSeq = append(chunkSeq,tokens[j+a])
		}
        chunk.Upsert(bson.M{"ngram":strings.Join(chunkSeq," ")},bson.M{"$inc": bson.M{"count": 1}})
	}
}

func InsertPOSNgram(tokens []string, n int) {
    maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C("posTags")
	for j := 0; j < (len(tokens) - n); j++ {
        posSeq := make([]string,0)
        otherSeq := make([]string,0)
        for a := 0; a < n; a++ {
            var temp = strings.Split(tokens[j+a],"/")
            if(len(temp) > 1) {
                posTag := temp[len(temp)-1]
                if(strings.Index(posTag,"+") != -1 && strings.Index(posTag,"+") != (len(posTag)-1)  ) {
                    temp = strings.Split(posTag,"+")
                    if(len(temp[1]) > 1 || temp[1] != "*") {
                        //pos.Upsert(bson.M{"posTag":temp[1]},bson.M{"$inc": bson.M{"count": 1}})
                        if (len(posSeq) == 1) {
                            pos.Upsert(&Ngrams{Ngram:posSeq[0]+" "+temp[1]},bson.M{"$inc": bson.M{"count": 1}})
                            if(len(otherSeq) == 1) {
                                fmt.Printf("%s----------",otherSeq[0]+otherSeq[1])
                                pos.Upsert(bson.M{"ngram":otherSeq[0]+" "+otherSeq[1]},bson.M{"$inc": bson.M{"count": 1}})
                            }
                        } else {
                            otherSeq = append(otherSeq,temp[1])
                        }
                    }
                    posTag = temp[0]
                }
                if(strings.Contains(posTag,"-")) {
                    posTag = strings.Split(posTag,"-")[0]
                }
                posSeq = append(posSeq,posTag)
            }
        }
        if (len(otherSeq) == 1 && len(posSeq) > 1) {
            // fmt.Printf("%v\n",tokens)
            // fmt.Printf("%s++++++++++\n",otherSeq[0]+posSeq[1])
            pos.Upsert(bson.M{"ngram":otherSeq[0]+" "+posSeq[1]},bson.M{"$inc": bson.M{"count": 1}})
        }
		if(len(posSeq) > 1) {
        	pos.Upsert(bson.M{"ngram":strings.Join(posSeq," ")},bson.M{"$inc": bson.M{"count": 1}})
		}
	}
}

func getNgram(nGram string, dbname string) int {
    session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C(dbname)
	results := Ngrams{}
	pos.Find(bson.M{"ngram": nGram}).One(&results)
	//fmt.Printf("%#v",results.Count)
	return results.Count
}

func getAllNgram(nGram []string, dbname string) []Ngrams {
    session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C(dbname)
	var results []Ngrams
	pos.Find(bson.M{"ngram": bson.M{"$in":nGram}}).All(&results)
	// fmt.Printf("%#v",posGram)
	// fmt.Printf("%v==========\n",len(results))
	return results
}

func getSomeNgram(start string, dbname string) Ngrams {
    session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C(dbname)
	results := Ngrams{}
	err = pos.Find(bson.M{"ngram": bson.M{"$regex":bson.RegEx{`^`+start+`.*`, ""}}}).Sort("-count").One(&results)
	// fmt.Printf("%v==========\n",len(results))
	return results
}

func getNSomeNgram(start string, dbname string) []Ngrams {
    session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C(dbname)
	var results []Ngrams
	err = pos.Find(bson.M{"ngram": bson.M{"$regex":bson.RegEx{`^`+start+`.*`, ""}}}).Sort("-count").All(&results)
	// fmt.Printf("%v==========\n",len(results))
	return results
}

func getWordPosgram(word string) []PosWordGram {
	maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C("wordPosgram")
	var results []PosWordGram
	pos.Find(bson.M{"word":word}).All(&results)
	// fmt.Printf("%#v",results)
	return results
}

func getAllWordPosgram(word []string) []PosWordGram {
	maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C("wordPosgram")
	var results []PosWordGram
	pos.Find(bson.M{"word":bson.M{"$in":word}}).All(&results)
	// fmt.Printf("%#v",results)
	return results
}

func getChunkPosgram(word string,postag string) []ChunkPosGram {
	maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C("chunkgrams")
	var results []ChunkPosGram
	pos.Find(bson.M{"word_pos":word+"-"+postag}).All(&results)
	return results
}

func getAllPosUnigrams() []PosUniGram {
	maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C("posUnigram")
	var results []PosUniGram
	pos.Find(bson.M{}).All(&results)
	// fmt.Printf("%#v",results)
	return results
}

func getAllChunkUnigrams() []Ngrams {
	maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C("chunkunigram")
	var results []Ngrams
	pos.Find(bson.M{}).All(&results)
	// fmt.Printf("%#v",results)
	return results
}

func getFrequencyAggregation(dbName string) map[int] int {
	maxWait := time.Duration(5 * time.Second)
    session, err := mgo.DialWithTimeout("127.0.0.1",maxWait)
	if err != nil {
		panic(err)
	}
    defer session.Close()
	session.SetMode(mgo.Monotonic, true)
    pos := session.DB("nlprokz").C(dbName)
	result := []bson.M{}
	pos.Pipe([]bson.M{bson.M{"$match": bson.M{}},bson.M{"$group":bson.M{"_id":"$count","count":bson.M{"$sum":1}}},bson.M{"$sort":bson.M{"count":-1}}}).All(&result)
	//fmt.Printf("%v",len(result))
	results := make(map[int] int,len(result))
	for _,i := range result {
		//fmt.Printf("%v\n",i["_id"])
		results[i["_id"].(int)] = i["count"].(int)
	}
	//fmt.Printf("%v",len(results))
	return results
}
// func main() {
//
//
//
//
// 	// // Drop Database
// 	// if IsDrop {
// 	// 	err = session.DB("nlprokz").DropDatabase()
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// }
//
// 	// Collection People
//
//
// 	// Index
// 	// index := mgo.Index{
// 	// 	Key:        []string{"name", "phone"},
// 	// 	Unique:     true,
// 	// 	DropDups:   true,
// 	// 	Background: true,
// 	// 	Sparse:     true,
// 	// }
//
// 	// err = c.EnsureIndex(index)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
//
// 	// Insert Datas
//
//
// 	// Query One
// 	result := Person{}
// 	err = c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Phone", result)
//
// 	// Query All
// 	var results []Person
// 	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
//
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Results All: ", results)
//
// 	// Update
// 	colQuerier := bson.M{"name": "Ale"}
// 	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}
// 	err = c.Update(colQuerier, change)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// Query All
// 	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
//
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Results All: ", results)
// }
