package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

type AttachmentStore struct {
	maxSize int64 // maximum file size in bytes
}

func NewAttachmentStore() *AttachmentStore {
	return &AttachmentStore{
		maxSize: 10 * 1024 * 1024, // 10MB default limit
	}
}

func (s *AttachmentStore) Store(file multipart.File, header *multipart.FileHeader, userID string) (*Attachment, error) {
	if header.Size > s.maxSize {
		return nil, fmt.Errorf("file too large: %d bytes (max %d)", header.Size, s.maxSize)
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Calculate hash
	hash := sha256.Sum256(content)
	hashString := hex.EncodeToString(hash[:])

	// Create attachment
	attachment := &Attachment{
		ID:         hashString,
		Name:       header.Filename,
		Type:       header.Header.Get("Content-Type"),
		Size:       header.Size,
		Hash:       hashString,
		Data:       content,
		UploadedBy: userID,
		UploadedAt: time.Now(),
	}

	return attachment, nil
}

// IsCompressibleType returns whether a file type should be compressed
func IsCompressibleType(mimeType string) bool {
	// List of mime types that are already compressed
	compressedTypes := map[string]bool{
		"image/jpeg":                   true,
		"image/png":                    true,
		"image/gif":                    true,
		"image/webp":                   true,
		"video/mp4":                    true,
		"video/mpeg":                   true,
		"audio/mpeg":                   true,
		"audio/mp4":                    true,
		"application/zip":              true,
		"application/x-gzip":           true,
		"application/x-rar-compressed": true,
		"application/x-7z-compressed":  true,
	}

	return !compressedTypes[mimeType]
}

// GetAttachmentData returns the attachment data
func (n *Node) GetAttachment(attachmentID string) (*Attachment, error) {
	if attachment, exists := n.Attachments[attachmentID]; exists {
		return &attachment, nil
	}
	return nil, fmt.Errorf("attachment not found: %s", attachmentID)
}

// GetEntryAttachment returns an attachment from an event entry
func (n *Node) GetEntryAttachment(eventID string, entryIndex int, attachmentID string) (*Attachment, error) {
	event, exists := n.Events[eventID]
	if !exists {
		return nil, fmt.Errorf("event not found: %s", eventID)
	}

	if entryIndex < 0 || entryIndex >= len(event.Entries) {
		return nil, fmt.Errorf("invalid entry index: %d", entryIndex)
	}

	for _, attachment := range event.Entries[entryIndex].Attachments {
		if attachment.ID == attachmentID {
			return &attachment, nil
		}
	}

	return nil, fmt.Errorf("attachment not found: %s", attachmentID)
}
