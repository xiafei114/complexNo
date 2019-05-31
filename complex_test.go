package seqno

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestNext(t *testing.T) {
	db, err := gorm.Open("mysql", "bindo:maxinfo@tcp(127.0.0.1:3306)/bindo?parseTime=true")
	if err != nil {
		fmt.Println(fmt.Sprintln(err))
	}

	defer db.Close()

	InitTable(db)

	no := NewNoGenerator(db, "Order")

	num, err := no.Next()
	if err != nil {
		fmt.Println(fmt.Sprintln(err))
	}
	fmt.Printf("%s\n", fmt.Sprintf(no.formateCode, "M", num))
}

// func TestGoroutine(t *testing.T) {
// 	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?parseTime=true")
// 	if err != nil {
// 		fmt.Println(fmt.Sprintln(err))
// 	}

// 	defer db.Close()

// 	InitTable(db)

// 	arrs := [5]string{"test-1", "test-2", "test-3", "test-4", "test-5"}

// 	var mutex sync.Mutex

// 	for _, v := range arrs {
// 		fmt.Printf("start %s\n", v)
// 		go func(v string) {
// 			for n := 0; n < 100; n++ {
// 				doSeqNoGenerator(db, v, &mutex)
// 			}
// 			time.Sleep(2 * time.Second)
// 			fmt.Printf("end %s\n", v)
// 		}(v)

// 	}

// 	time.Sleep(2 * time.Second)
// 	fmt.Println("end")
// }

// func doSeqNoGenerator(db *gorm.DB, logicID string, mutex *sync.Mutex) {
// 	mutex.Lock()
// 	num, err := NewSeqNoGenerator(db, logicID).Next()
// 	if err != nil {
// 		fmt.Println(fmt.Sprintln(err))
// 	}
// 	fmt.Printf("%s --- %s\n", logicID, num)
// 	mutex.Unlock()
// }
