import { Notification } from "./notification";

export class NotificationDOMFactory {
	protected readonly root_document: HTMLDocument;

	public constructor(root_document: HTMLDocument) {
		this.root_document = root_document;
	}

	public create_notification_dom(notification: Notification): HTMLElement {
		const wrapper = this.root_document.createElement('div');
		const title = this.root_document.createElement('span');
		const title_text = this.root_document.createTextNode(notification.title);
		title.append(title_text);
		wrapper.appendChild(title);
		wrapper.classList.add("notification");
		return wrapper;
	}
}
