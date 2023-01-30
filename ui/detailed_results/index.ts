import { Database } from '../core/proto_utils/database.js';
import { WindowedDetailedResults } from '../core/components/detailed_results.js';

Database.get();

const urlParams = new URLSearchParams(window.location.search);

if (urlParams.has('cssClass')) {
	document.body.classList.add(urlParams.get('cssClass')!);
}

const isIndividualSim = urlParams.has('isIndividualSim');
if (isIndividualSim) {
	document.body.classList.add('individual-sim');
}

document.body.classList.add('new-tab');

const detailedResults = new WindowedDetailedResults(document.body)
