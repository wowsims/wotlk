import { Database } from '../core/proto_utils/database.js';
import { WindowedDetailedResults } from '../core/components/detailed_results.js';

Database.get();

const urlParams = new URLSearchParams(window.location.search);

if (urlParams.has('themeBackgroundColor')) {
	document.body.style.setProperty('--theme-background-color', urlParams.get('themeBackgroundColor')!);
}
if (urlParams.has('themeBackgroundImage')) {
	document.body.style.setProperty('--theme-background-image', urlParams.get('themeBackgroundImage')!);
}
if (urlParams.has('themeBackgroundOpacity')) {
	document.body.style.setProperty('--theme-background-opacity', urlParams.get('themeBackgroundOpacity')!);
}

const isIndividualSim = urlParams.has('isIndividualSim');
if (isIndividualSim) {
	document.body.classList.add('individual-sim');
}

document.body.classList.add('new-tab');

const detailedResults = new WindowedDetailedResults(document.body)
