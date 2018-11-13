//    This file is part of tplwalker.
//
//    tplwalker is free software: you can redistribute it and/or modify it
//    under the terms of the GNU General Public License as published by the
//    Free Software Foundation, either version 3 of the License, or (at your
//    option) any later version.
//
//    tplwalker is distributed in the hope that it will be useful, but WITHOUT
//    ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
//    FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
//    more details.
//
//    You should have received a copy of the GNU General Public License along
//    with tplwalker.  If not, see <http://www.gnu.org/licenses/>.

// Package tplwalker is a tool to copy and transform a directory of templates
// and other files to a target directory.  It will re-create subdirectories,
// regular files and symlinks from a source directory.  Files that end with a
// configurable suffix, will be stripped of that suffix and be digested as a Go
// text template and executed with some data object.
//
// It was once intended as part of an engine that automatically creates
// directories of configuration files, but it has not been seen in action, yet.
package tplwalker
