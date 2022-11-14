const wowheadSupportedLanguages = [
	'en', // English
	'cn', // Chinese
	'de', // German
	'es', // Spanish
	'fr', // French
	'ko', // Korean
	'pt', // Portuguese
	'ru', // Russian
];

let cachedLanguageCode: string|null = null;

// Returns a 2-letter language code if it is a wowhead-supported language, or '' otherwise.
export function getLanguageCode(): string {
	if (cachedLanguageCode == null) {
		const browserLang = (navigator.language || '').substring(0, 2);
		if (wowheadSupportedLanguages.includes(browserLang)) {
			cachedLanguageCode = browserLang;
		} else {
			cachedLanguageCode = '';
		}
	}

	return cachedLanguageCode;
}
