package redup

import (
	"errors"
	"redup/internal/model"
)

func (ru *redupUsecase) AddBookmark(userID, videoID string) error {

	bookmarks, err := ru.bookmarkRepo.GetBookmarks(userID)
	if err != nil {
		return err
	}
	for _, bookmark := range bookmarks {
		if bookmark.VideoID == videoID {
			return errors.New("video already bookmarked")
		}
	}

	err = ru.bookmarkRepo.AddBookmark(userID, videoID)
	if err != nil {
		return err
	}
	return nil
}

func (ru *redupUsecase) RemoveBookmark(userID, bookmarkID string) error {

	err := ru.bookmarkRepo.RemoveBookmark(userID, bookmarkID)
	if err != nil {
		return err
	}
	return nil
}

func (ru *redupUsecase) GetBookmarks(userID string) ([]model.Bookmark, error) {

	bookmarks, err := ru.bookmarkRepo.GetBookmarks(userID)
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}
