#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

function filterExcludedFiles {
CHECK=$(git diff --name-only HEAD * \
		| grep -v .git \
		| grep -v ^vendor/ \
		| grep -v "\.txt$" \
		| grep -v "\.md$" \
		| grep -v "\.rst$" \
		| grep -v ^Gopkg\.lock$ | sort -u)
}
