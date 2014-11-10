/*
 * Copyright (c) 2014 Kurt Jung (Gmail: kurt.w.jung)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
Package recmgr provides a thin wrapper around Google's btree package that
facilitates the use of multiple indexes. This is particularly convenient for
managing structs that have multiple key fields.

This package operates on pointers to values, typically structs that can be
indexed in multiple ways.

The methods in this package correspond to the methods of the same name in the
btree package. Because multiple indexes are processed as a group, some methods
are not supported, for example DeleteMin() and DeleteMax(). Similarly, some
method semantics are different, for example Delete() returns the number of
removed keys rather than the deleted items.

Limitations

The records managed by this package are referenced by pointers so they should
remain accessible for the duration of the recmgr instance.

Non-key fields in these records can be changed with impunity. However, if key
fields are modified, it is advised to delete the record before modification and
add it again after to keep the underlying btrees consistent.
*/
package recmgr
