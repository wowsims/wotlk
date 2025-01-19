import { Tooltip } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element } from 'tsx-vanilla';

import { Component } from './component';

export class SocialLinks extends Component {
	static buildDiscordLink(): Element {
		const anchor = (
			<a
				href="https://discord.gg/p3DgvmnDCS"
				target="_blank"
				className="discord-link link-alt"
				dataset={{ bsToggle: 'tooltip', bsTitle: 'Join us on Discord' }}>
				<i className="fab fa-discord fa-lg" />
			</a>
		);
		Tooltip.getOrCreateInstance(anchor);
		return anchor;
	}

	static buildGitHubLink(): Element {
		const anchor = (
			<a
				href="https://github.com/wowsims/wotlk"
				target="_blank"
				className="github-link link-alt"
				dataset={{ bsToggle: 'tooltip', bsTitle: 'Contribute on GitHub' }}>
				<i className="fab fa-github fa-lg" />
			</a>
		);
		Tooltip.getOrCreateInstance(anchor);
		return anchor;
	}

	static buildPatreonLink(): Element {
		const anchor = (
			<a
				href="https://patreon.com/wowsims"
				target="_blank"
				className="patreon-link link-alt"
				dataset={{ bsToggle: 'tooltip', bsTitle: 'Support us on Patreon' }}>
				<i className="fab fa-patreon fa-lg" /> Patreon
			</a>
		);
		Tooltip.getOrCreateInstance(anchor);
		return anchor;
	}
}
