import { ActiveNotification } from "./active_notification";
import { NotificationLevel } from "./level";

export class Notification {
	public readonly title: string;
	public readonly level: NotificationLevel;
	public readonly duration: number;
	public readonly body?: string;

	public constructor(title: string, level: NotificationLevel, duration: number, body?: string) {
		this.title = title;
		this.level = level;
		this.duration = duration;
		this.body = body;
	}
}
