package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"pickup_dogshit/mlog"
	"time"
)

// @license.name MIT
// @title pickup_dogshit
// @author michaelliao
// @version 2019-02-08
// @version 0.1
// @description pick up dpogshit analysis
// @version 2024-09-24
// @version 0.2
// @description optimized code

// 扑克牌的13个数字(不记花色)
var ori_poker_nums = [13]byte{'A', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'J', 'Q', 'K'}

// 每次发牌前的所有扑克牌
var all_pokers []byte

func main() {
	//日志初始化
	mlog.InitLogger(&mlog.Params{
		Path:       "./log/pickup_dogshit.log", //文件路径
		MaxSize:    2,                          //MB 单个日志文件最大
		MaxBackups: 3,                          //备份个数
		MaxAge:     10,                         //保存时间,天
		Level:      -1,                         //# 日志级别
	})
	mlog.Debugf("pickup_dogshit start!")
	defer mlog.Debugf("pickup_dogshit end!")

	//读取参数:人数、牌的副数、循环次数
	ppeople_num := flag.Int("pp", 2, "people num")
	ppoker_num := flag.Int("pk", 1, "poker num")
	ploop_count := flag.Int("loop", 1000, "loop count")

	//
	flag.Usage = func() {
		fmt.Println("Usage: pickup_dpogshit -pp {people num} -pk {poker num} -loop {loop count}")
		os.Exit(1)
	}

	flag.Parse()
	people_num := *ppeople_num
	poker_num := *ppoker_num
	loop_count := *ploop_count

	mlog.Infof("people_num:%v, poker_num:%v, loop_count:%v", people_num, poker_num, loop_count)
	if people_num <= 0 || poker_num <= 0 {
		mlog.Errorf("people_num:%v, poker_num:%v", people_num, poker_num)
		return
	}

	all_poker_num := poker_num * 52
	all_pokers = make([]byte, all_poker_num)

	var all_num int64 = 0
	var max_num int = 0
	var min_num int = 100000
	//循环玩捡狗屎的游戏
	for i := 0; i < loop_count; i++ {
		loop_num := calc_dogshit(people_num, poker_num)
		all_num += int64(loop_num)

		if loop_num > max_num {
			max_num = loop_num
		}

		if loop_num < min_num {
			min_num = loop_num
		}
		mlog.Debugf("people_num:%v, poker_num:%v, loop_num:%v", people_num, poker_num, loop_num)
	}
	mlog.Infof("avg_num:%v, max_num:%v, min_num:%v", all_num/int64(loop_count), max_num, min_num)
}

func calc_dogshit(people_num int, poker_num int) int {
	// 使用当前时间作为种子初始化一个新的随机数源
	rs := rand.NewSource(time.Now().UnixNano())
	// 使用新的源创建一个新的随机数生成器
	rr := rand.New(rs)

	mlog.Debugf("calc_dogshit people_num:%v, poker_num:%v", people_num, poker_num)
	all_poker_num := poker_num * 52

	j := 0
	// 初始化poker
	for i := 0; i < all_poker_num; i++ {
		j = i % 13
		// mlog.Debugf("initial all_pokers j:%v, ori_poker_nums[j]:%v", j, ori_poker_nums[j])
		all_pokers[i] = ori_poker_nums[j]
	}

	mlog.Debugf("initial all_pokers:%v", all_pokers)

	// 发牌
	start_idxs := make([]int, people_num)
	end_idxs := make([]int, people_num)
	length_idxs := make([]int, people_num)
	var people_pokers [][]byte
	for i := 0; i < people_num; i++ {
		tmpArr := make([]byte, all_poker_num)
		people_pokers = append(people_pokers, tmpArr)
	}

	for i := 0; i < people_num; i++ {
		start_idxs[i] = 0
		end_idxs[i] = 0
	}

	k := 0
	for i := 0; i < all_poker_num; i++ {
		rand_num := rr.Intn(all_poker_num - i)
		count := 0
		for j = 0; j < all_poker_num; j++ {
			if all_pokers[j] != '\x00' {
				count++
			}

			if count == (rand_num + 1) {
				break
			}
		}

		k = i % people_num
		// mlog.Debugf("fa pai rand_num:%v, i:%v, j:%v, k:%v, end_idxs[k]:%v", rand_num, i, j, k, end_idxs[k])
		people_pokers[k][end_idxs[k]] = all_pokers[j]
		end_idxs[k] += 1
		all_pokers[j] = '\x00'
	}

	mlog.Debugf("all_pokers:%v ", all_pokers)

	mlog.Debugf("start_idxs:%v", start_idxs)
	mlog.Debugf("end_idxs:%v", end_idxs)
	for i := 0; i < people_num; i++ {
		mlog.Debugf("people idx:%v, pokers :%v", i, people_pokers[i])
	}

	// 打牌过程
	ds_curr_idx := -1
	loser_num := 0
	loop_num := 0
	match := 0
	for {
		loser_num = 0
		for i := 0; i < people_num; i++ {
			if start_idxs[i] == end_idxs[i] {
				loser_num++
				continue
			}

			poker_val := people_pokers[i][start_idxs[i]]
			match = 0
			for j = 0; j <= ds_curr_idx; j++ {
				if poker_val == all_pokers[j] {
					for k = j; k <= ds_curr_idx; k++ {
						people_pokers[i][end_idxs[i]] = all_pokers[k]
						all_pokers[k] = '\x00'
						end_idxs[i] = (end_idxs[i] + 1) % all_poker_num
					}

					people_pokers[i][end_idxs[i]] = poker_val
					end_idxs[i] = (end_idxs[i] + 1) % all_poker_num

					ds_curr_idx = j - 1
					match = 1
				}

			}

			people_pokers[i][start_idxs[i]] = '\x00'

			start_idxs[i] = (start_idxs[i] + 1) % all_poker_num
			if match == 0 {
				ds_curr_idx += 1
				all_pokers[ds_curr_idx] = poker_val
			}

			// 计算长度
			if start_idxs[i] > end_idxs[i] {
				length_idxs[i] = end_idxs[i] + (all_poker_num - start_idxs[i])
			} else if start_idxs[i] < end_idxs[i] {
				length_idxs[i] = end_idxs[i] - start_idxs[i]
			} else {
				length_idxs[i] = 0
			}
		}

		loop_num++

		//当失败者的数目是总人数减1时结束游戏
		if loser_num >= (people_num - 1) {
			mlog.Debugf("loop_num:%v", loop_num)
			mlog.Debugf("all_pokers:%v", all_pokers)
			mlog.Debugf("start_idxs:%v", start_idxs)
			mlog.Debugf("end_idxs:%v", end_idxs)
			mlog.Debugf("length_idxs:%v", length_idxs)
			for i := 0; i < people_num; i++ {
				mlog.Debugf("people idx:%v, pokers :%v", i, people_pokers[i])
			}
		}

		if loser_num >= (people_num - 1) {
			for i := 0; i < people_num; i++ {
				if start_idxs[i] == end_idxs[i] {
					mlog.Debugf("loser_idx:%v", i)
				}
			}

			break
		}
	}

	return loop_num
}
