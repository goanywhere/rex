/**
 * ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 * ----------------------------------------------------------------------
 *  Copyright © 2014 GoAnywhere Ltd. All Rights Reserved.
 * ----------------------------------------------------------------------*/

package web

import (
	"html/template"
	"path/filepath"

	"github.com/goanywhere/web/env"
)

// TODO boost ME.

var templates map[string]*template.Template

// ---------------------------------------------------------------------------
//  HTTP Response Rendering
// ---------------------------------------------------------------------------
// Forcely parse the passed in template files under the pre-defined template folder,
// & panics if the error is non-nil. It also try finding the default layout page (defined
// in ctx.Options.Layout) as the render base first, the parsed template page will be
// cached in global singleton holder.
func loadTemplates(filename string, others ...string) *template.Template {
	page, exists := templates[filename]
	if !exists {
		var files []string
		folder := env.Get("templates")
		files = append(files, filepath.Join(folder, filename))
		for _, item := range others {
			files = append(files, filepath.Join(folder, item))
		}

		page = template.Must(template.ParseFiles(files...))
		templates[filename] = page
	}
	return page
}

func init() {
	templates = make(map[string]*template.Template)
}
