export const CURRENT_PHASE = 5;

// Github pages serves our site under the /tbc directory (because the repo name is tbc)
export const REPO_NAME = 'tbc';

// Get 'elemental_shaman', the pathname part after the repo name
const pathnameParts = window.location.pathname.split('/');
const repoPartIdx = pathnameParts.findIndex(part => part == REPO_NAME);
export const SPEC_DIRECTORY = repoPartIdx == -1 ? '' : pathnameParts[repoPartIdx + 1];
