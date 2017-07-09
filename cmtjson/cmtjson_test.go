package cmtjson

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

type testData struct {
	Name   string
	Age    int
	Skills []string
	Cards  map[string]string
	IsMale bool `json:"is_male"`
}

var (
	testCases = map[string]string{
		// "test": `/* hahaha block comment hello! */`,
		"raw": `{
		"name":"iyidan",
		"Age":18,
		"skills":["go","python","php"],
		"Cards": {"ICBC":"123456","BC":"123455"},
		"is_male":true
	}`,
		"sharpcmt": ` #sharp cmt
		{ # sharp 你好
		"name":"iyidan", # sharp cmt
		"Age":18,
		#sharpcmt
		"skills":["go","python","php"],
		     # sharp cmt
		"Cards": {"ICBC":"123456","BC":"123455"},
		# "sharp" cmt {obj:1} // sharp /* sharp */
		"is_male":true,
		"xxx_Unrecognized": "\\#fawef\"\\#fadf",
		"xxx_Unrecognized2": "#fawef\"\\#fadf"
		# sharp
	}`,
		"slash": ` //slash cmt
		{ // slash 你好
		"name":"iyidan", // slash cmt
		"Age":18,
		//slash
		"skills":["go","python","php"],
		     // slash cmt
		"Cards": {"ICBC":"123456","BC":"123455"},
		// "slash" cmt {obj:1} # slash /* slash */
		"is_male":true,
		"xxx_Unrecognized": "//fawef\"//fadf", //o slash
		"xxx_Unrecognized2": "//fawef\"//#fadf" // slash cmt ##
		# slash
	}`,
		"blockcmt": ` /*blockcmt cmt*/
		{ /* blockcmt 你好 */
		"name":"iyidan", /* blockcmt cmt */
		"Age":18,/*blockcmt cmt */
		/*blockcmt
		* test
		* testddd
		*/
		"skills":["go","python","php"],
		/** blockcmt
		 * test
		 * testddd 
		 * @see http://test.com
		 * @author iyidan
		 */
		"Cards": {"ICBC":"123456","BC":"123455"},
		// "blockcmt" cmt {obj:1} # blockcmt /* blockcmt */
		"is_male":true,
		"xxx_Unrecognized": "//fawef\"//fadf", //o blockcmt
		"xxx_Unrecognized2": "//fawef\"//#fadf", // blockcmt cmt ##
		/* 
			"xxxtest":"test",
		*/
		"xxx_Unrecognized3": "/* test hello */", /* blockcmt cmt */ ## test
		"xxx_Unrecognized3": "/*test\" hello*/" /* blockcmt cmt */  // test
		/** blockcmt
		 * test
		 * testddd 
		 * @see http://test.com
		 * @author iyidan
		 */
		 # hh
	}`,
		"mixedcmt": `{ // comment
		/* hahaha block comment hello! */
		/** hahaha block comment hello! */
		// 你好，世界 // # "fawef" \n\tc
		# username ### // 你好
		# sharp cmt 2
		"name":"iyidan", // name 
		   #user address
		"Age":18, # age
		////// to many slash's "fd {"key":"value"}
		// slash cmt 2
		// slash cmt 3
		"skills":["go","python","php"],
		/* hahaha block comment hello! */
		"Cards": {"ICBC":"123456","BC":"123455"}, /* cards */
		"is_male":true,
		#dfasf // fasfdddf # /* afawe /* */
		"xxx_Unrecognized": "\\#fawef\"\\#fadf",
		"xxx_Unrecognized2": "http://tet.com/ab.df//\/\/df好",
		/** 
		 * tett1 // @ ##fafwe
		 * test2
		 */
		"xxx_Unrecognized3": "http://te/**t.com*//a/*b*/.df//\/\/df好"
	}`,
	}
)

func TestParseFromBytes(t *testing.T) {
	for caseName, cs := range testCases {
		dt := testData{}
		err := ParseFromBytes([]byte(cs), &dt)
		if err != nil {
			t.Fatalf("case: %s => parse error: %s\n", caseName, err)
		}
		if ok, correct := checkParseResult(&dt); !ok {
			t.Fatalf("case: %s => checkParseResult failed:\n\nparsed:\n%#v\ncorrect:\n%#v\n", caseName, dt, correct)
		}
		t.Logf("case: %s => ok, ret: %#v\n", caseName, dt)
	}
}

func BenchmarkParseFromBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, cs := range testCases {
			dt := testData{}
			err := ParseFromBytes([]byte(cs), &dt)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func TestParseFromFile(t *testing.T) {
	for caseName, cs := range testCases {
		tmpfilename, err := getTempfileWithJSON([]byte(cs))
		if err != nil {
			t.Fatal(err)
		}
		dt := testData{}
		err = ParseFromFile(tmpfilename, &dt)
		if err != nil {
			t.Fatalf("case: %s => parse error: %s\n", caseName, err)
		}
		if ok, correct := checkParseResult(&dt); !ok {
			t.Fatalf("case: %s => checkParseResult failed:\n\nparsed:\n%#v\ncorrect:\n%#v\n", caseName, dt, correct)
		}
		t.Logf("case: %s => ok, ret: %#v\n", caseName, dt)
	}
}

func BenchmarkParseFromFile(b *testing.B) {
	b.StopTimer()
	tmpfiles := make([]string, len(testCases))
	i := 0
	var err error
	for _, cs := range testCases {
		tmpfiles[i], err = getTempfileWithJSON([]byte(cs))
		if err != nil {
			b.Fatal(err)
		}
		i++
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, filename := range tmpfiles {
			dt := testData{}
			err = ParseFromFile(filename, &dt)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func TestParseFromLargeFile(t *testing.T) {
	padding := `# test cmt
	/* test cmt */
	// test  cmt
	/** test 
	 * test
	 * test
	 */`

	paddingcnts := 10 * WriteBufSize / len(padding)

	for caseName, cs := range testCases {
		data := []byte(cs)
		data[len(data)-1] = ' '
		buf := bytes.NewBuffer(data)
		for i := 0; i < paddingcnts; i++ {
			_, err := buf.WriteString(padding)
			if err != nil {
				t.Fatal(err)
			}
		}
		err := buf.WriteByte('}')
		if err != nil {
			t.Fatal(err)
		}
		tmpfilename, err := getTempfileWithJSON(buf.Bytes())
		t.Logf("large-tmpfile: %s\n", tmpfilename)
		if err != nil {
			t.Fatal(err)
		}
		dt := testData{}
		err = ParseFromFile(tmpfilename, &dt)
		if err != nil {
			t.Fatalf("case: %s => parse error: %s\n", caseName, err)
		}
		if ok, correct := checkParseResult(&dt); !ok {
			t.Fatalf("case: %s => checkParseResult failed:\n\nparsed:\n%#v\ncorrect:\n%#v\n", caseName, dt, correct)
		}
		t.Logf("case: %s => ok, ret: %#v\n", caseName, dt)
	}
}

func BenchmarkParseFromLargeFile(b *testing.B) {
	b.StopTimer()
	padding := `# test cmt
	/* test cmt */
	// test  cmt
	/** test 
	 * test
	 * test
	 */`

	paddingcnts := 10 * WriteBufSize / len(padding)

	tmpfiles := make([]string, len(testCases))
	i := 0
	var err error
	for _, cs := range testCases {
		data := []byte(cs)
		data[len(data)-1] = ' '
		buf := bytes.NewBuffer(data)
		for i := 0; i < paddingcnts; i++ {
			_, err := buf.WriteString(padding)
			if err != nil {
				b.Fatal(err)
			}
		}
		err := buf.WriteByte('}')
		if err != nil {
			b.Fatal(err)
		}
		tmpfiles[i], err = getTempfileWithJSON(buf.Bytes())
		if err != nil {
			b.Fatal(err)
		}
		i++
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, filename := range tmpfiles {
			dt := testData{}
			err = ParseFromFile(filename, &dt)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func getTempfileWithJSON(data []byte) (string, error) {
	tmpfile, err := ioutil.TempFile("", "getTempfileWithJSON")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(data)
	if err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}

func checkParseResult(parsed *testData) (bool, *testData) {
	correct := testData{}
	json.Unmarshal([]byte(testCases["raw"]), &correct)
	return reflect.DeepEqual(parsed, &correct), &correct
}
