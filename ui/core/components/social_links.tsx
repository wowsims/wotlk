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
				dataset={{ bsToggle: 'tooltip', bsTitle: '开发者Discord(需翻墙)' }}>
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
				dataset={{ bsToggle: 'tooltip', bsTitle: 'Github项目开源代码库' }}>
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

	static buildBilibiliLink(): Element {
		const anchor = (
			<a
				href="https://space.bilibili.com/919498"
				target="_blank"
				className="patreon-link link-alt"
				dataset={{ bsToggle: 'tooltip', bsTitle: '你可以通过B站<充电>或者新手盒子邀请码<财高八抖>来支持开发者! 如果你有任何纠错和意见也可以私信!' }}>
				<i className="fa-brands fa-bilibili" /> 支持用爱发电的开发者!
			</a>
		);
		Tooltip.getOrCreateInstance(anchor);
		return anchor;
	}
}
