package urlshort

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/alia5/urlshort/storinator"
)

/**
 * Note: This code is, currently, not very performant, nor scalable. (no fancy base62 shit.)
 * It is, however, more than enough for my own personal use.
 *
 * Also: No support for http urls, just https
 * Why? Just to annoy you! 😇
**/

// FUCK no constant arrays (slices, whatever...)
// var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
func alphabet() []rune {
	return []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

const defaultShortyLen = 6

// TODO: maybe make all the "constants" (*cough*) configurable.

var baseUrl = "https://localhost:7080"

func SetBaseUrl(url string) {
	baseUrl = url
}

func Shorten(url string) (shortUrl string, created bool, err error) {
	shortText := genRandomChars()
	shortUrl, created, err = ShortenWithName(url, shortText)
	if aeErr, ok := err.(*storinator.ShortTextExistsError); ok {
		// try again
		shortUrl, created, err = ShortenWithName(url, shortText)
		if err != nil {
			// TODO: proper.
			fmt.Println("Fuck this! I'm out!")
			errStr := fmt.Sprintf(
				`ShortText already exists ("%s" - "%s"), and we tried again, but still failed`,
				aeErr.ShortText,
				aeErr.FullLink)
			panic(errStr)
		}
	}
	return shortUrl, created, err
}

func ShortenWithName(url string, shortText string) (shorturl string, created bool, err error) {
	existingShort, created, err := storinator.CreateLink(shortText, stripHttp(url))
	if err != nil {
		return "", false, err
	}
	if !created {
		fmt.Printf("WARN: link %s already exists with \"%s\"\n", stripHttp(url), existingShort)
		shortText = existingShort
	}
	return stripDoubleSlash(fmt.Sprintf("%s/%s", baseUrl, shortText)), created, err
}

// also FUCK no optional params
// there may are suitable workarounds, but for bools
// another overload is too verbose, and misusing variadics feels like a dirty hack
type UnshortenParams struct {
	Text  string
	Click bool
}

func Unshorten(p UnshortenParams) (string, error) {
	shorty, err := storinator.FindLink(storinator.ShortText, p.Text, p.Click)
	if err != nil || shorty.FullLink == "" {
		return "", err
	}
	return fmt.Sprintf("https://%s", shorty.FullLink), err
}

func stripHttp(url string) string {
	regex := regexp.MustCompile(`^http(s)*?:\/\/`)
	return regex.ReplaceAllString(url, "")
}

func stripDoubleSlash(url string) string {
	regex := regexp.MustCompile(`\/*?`)
	return regex.ReplaceAllString(url, "")
}

func genRandomChars() string {
	seed := rand.NewSource(time.Now().UnixNano())
	srand := rand.New(seed)
	var shorty string
	alph := alphabet()
	for i := 0; i < defaultShortyLen; i++ {
		shorty += string(alph[srand.Intn(len(alph))])
	}
	return shorty
}
