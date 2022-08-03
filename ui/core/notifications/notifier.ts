import { NotificationDOMFactory } from "./dom_factory";
import { NotificationLevel } from "./level";
import { Notification } from "./notification";

export class Notifier {
	protected readonly wrapper: HTMLElement;
	protected readonly dom_factory: NotificationDOMFactory;

	public constructor(wrapper: HTMLElement, dom_factory: NotificationDOMFactory) {
		this.wrapper = wrapper;
		this.dom_factory = dom_factory;
	}

	protected write_notification_to_dom(notification: Notification): void {
		const node = this.dom_factory.create_notification_dom(notification);
		this.wrapper.appendChild(node);
		setTimeout(() => this.wrapper.removeChild(node), notification.duration);
	}

	public notify_debug(title: string, body?: string, duration = 10000): void {
		this.write_notification_to_dom(new Notification(title, NotificationLevel.DEBUG, duration, body))
	}

	public notify_info(title: string, body?: string, duration = 3000): void {
		this.write_notification_to_dom(new Notification(title, NotificationLevel.INFO, duration, body))
	}

	public notify_warning(title: string, body?: string, duration = 5000): void {
		this.write_notification_to_dom(new Notification(title, NotificationLevel.WARNING, duration, body))
	}

	public notify_error(title: string, body?: string, duration = 5000): void {
		this.write_notification_to_dom(new Notification(title, NotificationLevel.WARNING, duration, body))
	}
}
