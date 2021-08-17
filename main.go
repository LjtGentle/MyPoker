/*
扑克牌52张，花色黑桃spades，红心hearts，方块diamonds，草花clubs各13张，2-10，J，Q，K，A

Face：即2-10，J，Q，K，A，其中10用T来表示。

Color：即S(spades)、H(hearts)、D(diamonds)、C(clubs)

用 Face字母+小写Color字母表示一张牌，比如As表示黑桃A，其中A为牌面，s为spades，即黑桃，Ah即红心A，以此类推。
现分别给定任意两手牌各有5张，例如：AsAhQsQhQc，vs KsKhKdKc2c，请按德州扑克的大小规则来判断双方大小。

代码要求有通用性，可以任意输入一手牌或几手牌来进行比较。


5张
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Poker struct {
	Alice string `json:"alice"`
	Bob string	`json:"bob"`
	Result int	`json:"result"`
}

type Match struct {
	Matches []Poker `json:"matches"`
}

// 定义一个牌型
// 1.皇家同花顺
// 2.同花顺
// 3.四条
// 4.满堂彩（葫芦，三带二）
// 5.同花
// 6.顺子
// 7.三条
// 8.两对
// 9.一对
// 10.单张大牌
type CardType int

// 把数据读取出来 分别放在切片中
func ReadFile(filename string) (alices,bobs []string, results []int){
	buf ,err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var matches Match
	err = json.Unmarshal(buf,&matches)
	if err != nil {
		panic(err)
	}

	alices = make([]string,1024)
	bobs = make([]string,1024)
	results = make([]int,1024)

	for k, v := range matches.Matches {
		//fmt.Printf("k=%#v,v=%#v\n",k,v)
		alices[k] = v.Alice
		bobs[k] = v.Bob
		results[k] = v.Result
	}
	return
}

// 根据手牌判断 牌型
func JudgmentGroup(cards string)(cardType CardType,cardSizes, cardColors []string ) {
	// 遍历分开手牌的大小和颜色
	cardSizes = make([]string,5)
	cardColors = make([]string,5)
	for i,v := range cards {
		// fmt.Println("v=",string(v))
		if i%2 !=0 {
			cardColors[i/2] = string(v)
		}else {
			cardSizes[i/2] = string(v)
		}
	}

	// 先判断是否是同花
	var s, h, d,c int
	for _, v := range cardColors {

		switch v {
		case "h": // 红心
			h++
		case "s": // 黑桃
			s++
		case "d": // 方块
			d++
		case "c": // 草花
			c++
		default:
			fmt.Println("无法解析的花色")
		}
	}
	// 同花
	if s == 5 || h == 5 || d == 5 || c == 5 {
		// 同花后判断是不是顺子

		// 判断是不是顺子
		isShun := IsShunZi(cardSizes)

		if isShun == true {
			// 是 就同花顺
			cardType = 1
		}else {
			// 否 就同花
			cardType = 5
		}
		return

	}

	// 不是同花
	// 先判断是不是顺子
	isShun := IsShunZi(cardSizes)
	if isShun == true {
		// 顺子出函数
		cardType = 6
		return
	}
	//  判断重复张数
	// var sum []int
	sums := make([]int,len(cardSizes)/2 +1)
	for i := 0; i<len(cardSizes)/2 +1 ; i++ {
		for _,v := range cardSizes {
			fmt.Printf("v111=%#v--->",v)
			if v == cardSizes[i] {
				sums[i] ++
			}
		}
		fmt.Println("sum=",sums[i])
	}
	k1 := 0 // 记录三个
	k2 := 0 // 记录两个
	for _,v := range sums {
		if v >=4 {
			// 为四条出函数
			cardType = 3
			return
		}
		if v == 3 {
			k1++
		}
		if v == 2 {
			k2++
		}
	}
	fmt.Printf("k1=%#v,k2=%#v\n",k1,k2)
	// 三带二 出函数
	if k1 >=1 && k2 >=1 {
		cardType = 4
		return
	}
	// 三条
	if k1 >= 1 {
		cardType = 7
		return
	}
	// 两对
	if k2 >= 2 {
		cardType = 8
		return
	}
	// 一对
	for i:=len(cardSizes)/2 +1 ; i< len(cardSizes);i++ {
		for _,v := range cardSizes {
			fmt.Printf("v111=%#v--->",v)
			if v == cardSizes[i] {
				sums[i/2] ++
			}
		}
		fmt.Println("sum=",sums[i/2])
	}
	for _,v := range sums {
		if v >2 {
			cardType = 9
			return
		}
	}
	// 单张大牌
	cardType=10
	return
}

// 判断是不是顺子
func IsShunZi(cardSize []string)(shunZi bool){
	shunZi = false
	saves := make([]string,13)
	for _, v := range cardSize {
		switch v {
		case "2":
			saves[0] = v
		case "3":
			saves[1]= v
		case "4":
			saves[2] = v
		case "5":
			saves[3] = v
		case "6":
			saves[4] =v
		case "7":
			saves[5]=v
		case "8":
			saves[6]=v
		case "9":
			saves[7] =v
		case "T":
			saves[8] = v
		case "J":
			saves[9]=v
		case "Q":
			saves[10]=v
		case "K":
			saves[11]=v
		case "A":
			saves[12] = v
		default:
			fmt.Println("无法解析的扑克牌")
		}
	}
	// 判断数组是否连续
	sum := 0
	for _ ,v := range saves {
		if v != "" {
			sum ++
		}else {
			sum = 0
		}
		if sum >=5 {
			// break
			shunZi = true
			return
		}
	}
	return
}

// 下面是同类型的比较 ---------
func TranNum(card uint8) (f float64){
	switch  {

	}
}
func PokerTranNum(cardLen int,card string)(fs []float64) {

	fs = make([]float64,cardLen)
	for i:=0 ; i <len(card); i +=2 {
		TranNum(card[i])
	}
	// for i, v := range card {
	//
	// }
}
// 单张大牌的比较  -----
func SingleCardCompareSize(cardLen int ,cards ...string)() {

	sizess  :=make(map[int][]int32)
	colorss :=make(map[int][]int32)
	//oder := make(map[int]map[int][]int32)
	for i, v := range cards {
		fmt.Printf("i=%#v,v=%#v\n",i,v)
		sizess[i] = make([]int32,cardLen)
		colorss[i] = make([]int32,cardLen)
		for j, value := range v {
			// 分开 牌大小 颜色
			if j%2 !=0 {
				colorss[i][j/2] = value
			}else {
				sizess[i][j/2] = value
			}

		}
		// 排序

		//fmt.Printf("%#v----coloress:%#v,sizess:%#v\n",i,colorss[i],sizess[i])

	}



		// for i :=0; i<len(v); i+=2{
		// 	switch v[i] {
		// 	case 50:
		// 		seqs[0] = v[i]
		// 	case 51:
		// 		seqs[1] = value
		// 	case 52:
		// 		seqs[2] = value
		// 	case 53:
		// 		seqs[3] = value
		// 	case 54:
		// 		seqs[4] = value
		// 	case 55:
		// 		seqs[5] = value
		// 	case 56:
		// 		seqs[6] = value
		// 	case 57:
		// 		seqs[7] = value
		// 	case 84:
		// 		seqs[8] = value
		// 	case 74:
		// 		seqs[9] = value
		// 	case 81:
		// 		seqs[10]= value
		// 	case 75:
		// 		seqs[11]= value
		// 	case 65:
		// 		seqs[12]= value
		// 	}
		// }



}


func main() {
	// file := "/home/weilijie/chromeDown/match_result.json"
	// ReadFile(file)

	// JudgmentGroup("AsAhQsQhQc")
	str1 := "23456789"
	fmt.Println("[]byte(str1)",[]byte(str1))
	str2 := "T"
	fmt.Println("[]byte(str2)",[]byte(str2))
	str3 := "JQKA"
	fmt.Println("[]byte(str3)",[]byte(str3))
	str := []string {"2"}
	isres := IsShunZi(str)
	fmt.Println("isres=",isres)

	type1 ,_,_ := JudgmentGroup("KsKhKdKc2c")
	fmt.Println("type1=",type1)
	SingleCardCompareSize(5,"KsKhKdKc2c","AsAhQsQhQc")
	str4 := "shcd"
	fmt.Println("[]byte(str4)=",[]byte(str4))

	var int32 int32
	int32 = 105*1000
	fmt.Println("int_32=",int32)
}
