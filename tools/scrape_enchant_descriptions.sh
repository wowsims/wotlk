#!/usr/bin/env sh

# Usage example:
# ./tools/scrape_enchant_descriptions <./sim/core/items/all_enchants.go >./assets/enchants/descriptions.json

set -e

spell_url_from_item_id() {
	safe_curl "https://www.wowhead.com/wotlk/item=$1" |
		grep -o '\/wotlk\/spell=[0-9]\+' |
		head -n 1
}

id_from_entry() {
	grep -oP '{ID: \K[0-9]+'
}

delete_newlines() {
	tr -d '\n'
}

safe_curl() {
	curl -sL "$1" |
		remove_user_posts |
		delete_newlines
}

spell_effect_from_url() {
	safe_curl "https://www.wowhead.com$1" |
		grep -oP '<th>Effect.*?</th>.*<a href="\/wotlk\/spell=[0-9]+/.*?">\K.*?(?=</a>)|<th>Effect.*?</th>.*<span class="q2">\K.*(?=</span>)'
}

remove_user_posts() {
	grep -v '<div class="user-post-">'
}

printf '{'
grep '^\s*{' |
	while read -r line; do
		id=$(echo "$line" | id_from_entry)
		effect=""
		if echo "$line" | grep 'IsSpellID: true' >/dev/null; then
			effect=$(spell_effect_from_url "/wotlk/spell=$id")
		else
			spell_url_suffix=$(spell_url_from_item_id "$id")
			effect=$(spell_effect_from_url "$spell_url_suffix")
		fi

		printf '\n  "%d": "%s",' "$id" "$effect"
		sleep 1 # Avoid potential rate-limting
	done |
	sed '$s/\,//g'
printf '\n}\n'
