package models

import (
	"encoding/json"
	"testing"
)

func TestArticleInitialization(t *testing.T) {
	article := Article{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
	}

	if article.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", article.ID)
	}
	if article.Title != "Test Title" {
		t.Errorf("expected Title to be 'Test Title', got '%s'", article.Title)
	}
	if article.Content != "Test Content" {
		t.Errorf("expected Content to be 'Test Content', got '%s'", article.Content)
	}
}

func TestArticleSerialization(t *testing.T) {
	article := Article{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
	}

	data, err := json.Marshal(article)
	if err != nil {
		t.Fatalf("failed to serialize article: %v", err)
	}

	expectedJSON := `{"id":1,"title":"Test Title","content":"Test Content"}`
	if string(data) != expectedJSON {
		t.Errorf("expected JSON '%s', got '%s'", expectedJSON, string(data))
	}
}

func TestArticleDeserialization(t *testing.T) {
	data := `{"id":1,"title":"Test Title","content":"Test Content"}`
	var article Article

	err := json.Unmarshal([]byte(data), &article)
	if err != nil {
		t.Fatalf("failed to deserialize JSON: %v", err)
	}

	if article.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", article.ID)
	}
	if article.Title != "Test Title" {
		t.Errorf("expected Title to be 'Test Title', got '%s'", article.Title)
	}
	if article.Content != "Test Content" {
		t.Errorf("expected Content to be 'Test Content', got '%s'", article.Content)
	}
}
