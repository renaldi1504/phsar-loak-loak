package seller

import "testing"

// Video opsional: produk foto-saja boleh, asal fotonya cukup (minPhotosNoVideo).
func TestCheckSubmissionVideoOptional(t *testing.T) {
	media := func(kinds ...string) []MediaInput {
		out := make([]MediaInput, len(kinds))
		for i, k := range kinds {
			out[i] = MediaInput{Key: "k", Kind: k}
		}
		return out
	}
	cases := []struct {
		name  string
		cond  string
		kinds []string
		ok    bool
	}{
		{"video + 1 foto", "new", []string{"video", "photo"}, true},
		{"3 foto tanpa video", "new", []string{"photo", "photo", "photo"}, true},
		{"2 foto tanpa video", "new", []string{"photo", "photo"}, false},
		{"video saja", "new", []string{"video"}, false},
		{"kosong", "new", nil, false},
		{"bekas tanpa foto cacat", "used_good", []string{"photo", "photo", "photo"}, false},
		{"bekas lengkap tanpa video", "used_good", []string{"photo", "photo", "photo", "defect_photo"}, true},
	}
	for _, c := range cases {
		err := checkSubmission(ProductInput{Condition: c.cond, Media: media(c.kinds...)})
		if (err == nil) != c.ok {
			t.Errorf("%s: err=%v, want ok=%v", c.name, err, c.ok)
		}
	}
}
