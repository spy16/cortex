# Cortex

Cortex is the backend for the chunked app. It maintains chunks and provides API for managing them.

## Design

### Chunk

Everything including (but not limited to) notes, todo-list, reminders, etc. are all chunks. A chunk will have the
following structure:

```golang
package chunk

import "time"

type Chunk struct {
	ID        string    // A unique 1-12 character string identifier for the chunk
	Type      string    // One of reminder, text, image, snippet, etc.
	Data      string    // content of the chunk formatted as per type.
	Score     int64     // A priority score to order neighboring chunks.
	Parent    string    // identifier for the parent chunk.
	Author    string    // ID of the user who created it.
	CreatedAt time.Time // timestamp at which the chunk was added.
	UpdatedAt time.Time // timestamp at which the chunk was updated.
}
```

#### A **Note Chunk** without any parent:

 ```json
 {
  "id": "4a35aed",
  "type": "NOTE",
  "data": {
    "format": "markdown",
    "text": "This is a sample note chunk. A note will contain rich-text in some format indicated by the 'format' key."
  },
  "author": "spy16",
  "created_at": "2022-05-10T10:30:00.000Z",
  "updated_at": "2022-05-10T10:30:00.000Z"
}
 ```

#### A **Note Chunk** with a parent:

```json
{
  "id": "4a25aec",
  "type": "NOTE",
  "parent": "4a35aed",
  "data": {
    "format": "markdown",
    "text": "This is a sample note chunk. A note will contain rich-text in some format indicated by the 'format' key."
  },
  "author": "spy16",
  "created_at": "2022-05-10T10:30:00.000Z",
  "updated_at": "2022-05-10T10:30:00.000Z"
}
```

#### An **Image chunk**:

```json
{
  "id": "4a25aec",
  "type": "IMAGE",
  "parent": "4a35aed",
  "data": {
    "caption": "My cute little dog",
    "alt": "markdown",
    "url": "https://<url-to-the-image>"
  },
  "author": "spy16",
  "created_at": "2022-05-10T10:30:00.000Z",
  "updated_at": "2022-05-10T10:30:00.000Z"
}
```

#### A **todo chunk**:

```json
{
  "id": "4a25aec",
  "type": "TODO",
  "parent": "4a35aed",
  "data": {
    "deadline": "",
    "items": [
      {
        "text": "Get groceries"
      },
      {
        "text": "Pay the bills"
      }
    ]
  },
  "author": "spy16",
  "created_at": "2022-05-10T10:30:00.000Z",
  "updated_at": "2022-05-10T10:30:00.000Z"
}
```