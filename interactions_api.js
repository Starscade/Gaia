import { GoogleGenAI } from 'https://esm.sh/@google/genai'

export default class {
	constructor({
		api_key,
		intellect,
		model,
		persona,
		print_function = (text) => {
			console.log(text)
		},
	} = {}) {
		this.API_KEY = api_key ?? ''
		this.INTELLECT = intellect ?? 'low'
		this.MODEL = model ?? 'gemini-flash-lite-latest'
		this.PERSONA = persona ?? 'Your name is Gaia. Respond in short, plaintext SMS with the occasional emoji.'
		this.STDOUT = print_function
		this.TRANSCRIPT_STORAGE_KEY = 'GAIA_TOPIC_ID'

		this.AI = new GoogleGenAI({
			apiKey: this.API_KEY,
		})
	}

	async ask({
		attachments = [],
		cmd = '',
		preserve_context = false,
		modalities = [
			'text',
		],
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
			attachments.forEach( attachment => {
				interaction_obj.input.push({
					data: attachment.data,
					mime_type: attachment.mime_type,
					resolution: this.INTELLECT,
					type: attachment.mime_type.split('/')[0],
				})
			})
		}

		// console.debug(interaction_obj)

		const interaction = await this.AI.interactions.create(interaction_obj)

		// console.debug(interaction)

		for await (const event of interaction) {
			// console.debug(event)
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

	}

	forgetTranscript() {
		const topic_id = localStorage.removeItem(this.TRANSCRIPT_STORAGE_KEY)
		return null
	}

	async getEcho() {
		const transcript = await this.getTranscript()
		return transcript?.output_image.data ?? transcript.output_text
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

