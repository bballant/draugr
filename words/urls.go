package words

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var counts []string = []string{
	"900,google.com",
	"60,mail.yahoo.com",
	"10,mobile.sports.yahoo.com",
	"40,sports.yahoo.com",
	"300,yahoo.com",
	"10,stackoverflow.com",
	"2,en.wikipedia.org",
	"1,es.wikipedia.org",
}

type subsum struct {
	first string
	count int
}

func headTail[A comparable](l []A) (A, []A) {
	if len(l) == 0 {
		panic("empty slice")
	}
	if len(l) == 1 {
		return l[0], nil
	}
	return l[0], l[1:]
}

func tracePath(acc []string, p1 []string, p2 []string) []string {
	if len(p1) == 0 || len(p2) == 0 {
		return acc
	}
	h1, t1 := headTail(p1)
	h2, t2 := headTail(p2)
	if h1 != h2 {
		return acc
	}
	return tracePath(append(acc, h1), t1, t2)
}

func findLongestPathBetween(p1 []string, p2 []string) []string {
	var innerLongest func([]string, []string, []string) []string
	innerLongest = func(longest []string, _p1 []string, _p2 []string) []string {
		if _p1 == nil {
			return longest
		}
		if _p2 == nil {
			_, _t1 := headTail(_p1)
			return innerLongest(longest, _t1, p2)
		}
		path := tracePath([]string{}, _p1, _p2)
		_, _t2 := headTail(_p2)
		if len(path) > len(longest) {
			return innerLongest(path, _p1, _t2)
		}
		return innerLongest(longest, _p1, _t2)
	}
	return innerLongest([]string{}, p1, p2)
}

func gimmeCounts(cs []string) []subsum {
	var innerCounts func(acc map[string]int, cs_ []string) map[string]int
	innerCounts = func(acc map[string]int, cs_ []string) map[string]int {
		if len(cs_) == 0 {
			return acc
		}
		h := cs_[0]
		row := strings.Split(h, ",")
		if len(row) != 2 {
			panic("woah!")
		}
		num, err := strconv.Atoi(row[0])
		if err != nil {
			panic("can't convert string to integer " + row[0])
		}
		s := strings.Split(row[1], ".")
		for _, v := range s {
			if _, ok := acc[v]; !ok {
				acc[v] = num
			} else {
				acc[v] = acc[v] + num
			}
		}
		return innerCounts(acc, cs_[1:])
	}
	resMap := innerCounts(map[string]int{}, cs)
	resList := []subsum{}
	for k, v := range resMap {
		resList = append(resList, subsum{k, v})
	}
	sort.SliceStable(resList, func(i, j int) bool {
		return resList[i].count < resList[j].count
	})
	return resList
}

func UrlsRun() {
	fmt.Println(gimmeCounts(counts))
	user0 := []string{"/nine.html", "/four.html", "/six.html", "/seven.html", "/one.html"}
	user2 := []string{"/nine.html", "/two.html", "/three.html", "/four.html", "/six.html",
		"/seven.html"}
	fmt.Println(findLongestPathBetween(user0, user2))
}

/*
 *
 *
 ;; counts = [ "900,google.com",
;;      "60,mail.yahoo.com",
;;      "10,mobile.sports.yahoo.com",
;;      "40,sports.yahoo.com",
;;      "300,yahoo.com",
;;      "10,stackoverflow.com",
;;      "2,en.wikipedia.org",
;;      "1,es.wikipedia.org" ]

;; Expected output (in any order, any format):
;; 1320    com
;;  900    google.com
;;  410    yahoo.com
;;   60    mail.yahoo.com
;;   10    mobile.sports.yahoo.com
;;   50    sports.yahoo.com
;;   10    stackoverflow.com
;;    3    org
;;    3    wikipedia.org
;;    2    en.wikipedia.org
;;    1    es.wikipedia.org

;; You are in charge of a display advertising program. Your ads are displayed on websites all over the internet. You have some CSV input data that counts how many times you showed an ad on each individual domain. Every line consists of a count and a domain name.

;; Write a function to return the total number of hits for each domain and subdomain:
;;     param - counts
;;     return -  a data structure containing the number of hits that were recorded on each domain AND each domain under it

;; For example, a click on "sports.yahoo.com" counts for "sports.yahoo.com", "yahoo.com", and "com". (Subdomains are added to the left of their parent domain. So "sports" and "sports.yahoo" are not valid domains.)


;;   String[] counts = {
;;       "900,google.com",
;;       "60,mail.yahoo.com",
;;       "10,mobile.sports.yahoo.com",
;;       "40,sports.yahoo.com",
;;       "300,yahoo.com",
;;       "10,stackoverflow.com",
;;       "2,en.wikipedia.org",
;;       "1,es.wikipedia.org" };

;; We have some clickstream data that we gathered on our client's website. Using cookies, we collected snippets of users' anonymized URL histories while they browsed the site. The histories are in chronological order and no URL was visited more than once per person.

;; Write a function that takes two user’s browsing histories as input and returns the longest contiguous sequence of URLs that appears in both.

;; Sample input:

;; user0 = [ "/nine.html", "/four.html", "/six.html", "/seven.html", "/one.html" ]
;; user1 = [ "/one.html", "/two.html", "/three.html", "/four.html", "/six.html" ]
;; user2 = [ "/nine.html", "/two.html", "/three.html", "/four.html", "/six.html",
;;           "/seven.html" ]
;; user3 = [ "/three.html", "/eight.html" ]

;; (user0, user2):
;;    /four.html
;;    /six.html
;;    /seven.html*
 *
 *
 *
*/
