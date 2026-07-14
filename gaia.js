import { GoogleGenAI } from '@gemini'

export default class {
	constructor({
		api_key = '',
		debug = false,
		intellect = 'low',
		model = 'gemini-flash-lite-latest',
		persona =
			'Your name is Gaia. Respond in short, plaintext SMS with the occasional emoji.',
		print_function = (text) => {
			console.log(text)
		},
	} = {}) {
		this.API_KEY = api_key
		this.DEBUG = debug
		this.INTELLECT = intellect
		this.MODEL = model
		this.PERSONA = persona
		this.STDOUT = print_function
		this.TRANSCRIPT_STORAGE_KEY = 'GAIA_TOPIC_ID'

		this.AI = new GoogleGenAI({
			apiKey: this.API_KEY,
		})
	}

	async ask({
		attachments = [],
		cmd = '',
		modalities = [
			'text',
		],
		preserve_context = false,
		user_prompt = '',
	} = {}) {
		if (cmd) {
			switch (cmd.trim().toUpperCase()) {
				case 'ECHO':
					this.STDOUT(await this.getEcho())
					break
			}
		}

		if (!user_prompt && attachments.length === 0) {
			return {
				err: 'No prompt!',
			}
		}

		const interaction_obj = {
			input: [
				{
					text: user_prompt,
					type: 'text',
				},
			],
			generation_config: {
				image_config: {
					image_size: '1K',
				},
				thinking_level: this.INTELLECT,
			},
			model: this.MODEL,
			response_modalities: modalities,
			stream: true,
			system_instruction: this.PERSONA,
		}

		if (modalities.includes('image')) {
			interaction_obj.model = 'gemini-3.1-flash-lite-image'
		}

		if (preserve_context) {
			const topic_id = localStorage.getItem(this.TRANSCRIPT_STORAGE_KEY)
			if (topic_id) {
				interaction_obj.previous_interaction_id = topic_id
			}
		}

		if (attachments.length > 0) {
			attachments.forEach((attachment) => {
				const template_attachment = {
					type: attachment.mime_type.split('/')[0],
				}

				if (attachment.mime_type === 'text/plain') {
					template_attachment.text = attachment.data
				} else {
					template_attachment.data = attachment.data
					template_attachment.mime_type = attachment.mime_type
					template_attachment.resolution = this.INTELLECT
				}

				interaction_obj.input.push(template_attachment)
			})
		}

		if (this.DEBUG) {
			console.debug({
				debug: interaction_obj,
				name: 'INTERACTION_OBJECT',
			})
		}

		let interaction

		try {
			interaction = await this.AI.interactions.create(interaction_obj)
		} catch (err) {
			return {
				err: err.error.error.message,
			}
		}

		if (this.DEBUG) {
			console.debug({
				debug: interaction,
				name: 'INTERACTION',
			})
		}

		for await (const event of interaction) {
			if (this.DEBUG) {
				console.debug({
					debug: event,
					name: 'EVENT',
				})
			}

			switch (event.event_type) {
				case 'error':
					this.STDOUT(event.error.message)
					break
				case 'interaction.created':
					localStorage.setItem('GAIA_TOPIC_ID', event.interaction.id)
					break
				case 'step.delta':
					switch (event.delta.type) {
						case 'text':
							this.STDOUT(event.delta.text)
							break
						default:
							if (event.delta.data) {
								this.STDOUT(event.delta.data)
							}
							break
					}
					break
			}
		}

		return {
			err: null,
		}
	}

	forgetTranscript() {
		localStorage.removeItem(this.TRANSCRIPT_STORAGE_KEY)
		return null
	}

	async getEcho() {
		const transcript = await this.getTranscript()
		return transcript?.output_image?.data ?? transcript.output_text
	}

	getEnv() {
		return {
			GAIA_API_KEY: this.API_KEY,
			GAIA_PERSONA: this.PERSONA,
		}
	}

	async getTranscript() {
		const topic_id = localStorage.getItem(this.TRANSCRIPT_STORAGE_KEY)
		const prior_interaction = await this.AI.interactions.get(topic_id)
		return prior_interaction
	}
}
