import { GoogleGenAI } from '@gemini'

export default class {
	constructor({
		api_key = '',
		debug = false,
		intellect = 'low',
		persona =
			'Your name is Gaia. Respond in short, plaintext SMS with the occasional emoji.',
		print_function = (model_output) => {
			console.log(model_output.data)
		},
	} = {}) {
		this.API_KEY = api_key
		this.DEBUG = debug
		this.INTELLECT = intellect
		this.MODELS = {
			audio: 'gemini-3.1-flash-tts-preview',
			image: 'gemini-3.1-flash-image',
			text: 'gemini-flash-lite-latest',
		}
		this.PERSONA = persona
		this.STDOUT = print_function
		this.TOPIC_STORAGE_KEY = 'GAIA_TOPIC_ID'
		this.TRANSCRIPT_STORAGE_KEY = 'GAIA_TRANSCRIPT'

		this.AI = new GoogleGenAI({
			apiKey: this.API_KEY,
		})

		this.getEnv()
	}

	async ask({
		attachments = [],
		cmd = '',
		modality = 'text',
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
				speech_config: [
					{ voice: 'Sulafat' },
				],
				thinking_level: this.INTELLECT,
			},
			model: this.MODELS[modality],
			response_modalities: [modality],
			stream: true,
			system_instruction: this.PERSONA,
		}

		if (preserve_context) {
			const topic_id = localStorage.getItem(this.TOPIC_STORAGE_KEY)
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

		if (modality === 'audio') {
			delete interaction_obj.generation_config.thinking_level
			delete interaction_obj.previous_interaction_id
			delete interaction_obj.system_instruction
		}

		this.printDebug('INTERACTION_OBJ', interaction_obj)

		let interaction

		try {
			interaction = await this.AI.interactions.create(interaction_obj)
		} catch (err) {
			this.printDebug('INTERACTION_ERR', err, true)
			return {
				err: err?.error?.error?.message ?? err,
			}
		}

		this.printDebug('INTERACTION', interaction)

		for await (const event of interaction) {
			this.printDebug('EVENT', event)

			switch (event.event_type) {
				case 'error':
					return {
						err: event.error.message,
					}
				case 'interaction.created':
					localStorage.setItem('GAIA_TOPIC_ID', event.interaction.id)
					break
				case 'step.delta':
					switch (event.delta.type) {
						case 'text':
							this.STDOUT({
								data: event.delta.text,
								type: event.delta.type,
							})
							break
						default:
							if (event.delta.data) {
								this.STDOUT({
									data: event.delta.data,
									type: event.delta.type,
								})
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
		localStorage.removeItem(this.TOPIC_STORAGE_KEY)
		return null
	}

	async getEcho() {
		const transcript = await this.getTranscript()
		if (transcript.err) {
			return {
				err: transcript.err,
			}
		}
		if (transcript) {
			const data_type = transcript.output_image
				? 'image'
				: (transcript.output_audio
					? 'audio'
					: (transcript.output_text ? 'text' : null))

			return {
				data: transcript?.output_image?.data ??
					transcript?.output_audio?.data ??
					transcript?.output_text,
				type: data_type,
			}
		}
		return ''
	}

	getEnv() {
		const env_obj = {
			GAIA_API_KEY: this.API_KEY,
			GAIA_PERSONA: this.PERSONA,
		}
		this.printDebug('ENVIRONMENT', env_obj)
		return env_obj
	}

	async getTranscript() {
		const topic_id = localStorage.getItem(this.TOPIC_STORAGE_KEY)
		if (topic_id) {
			let prior_interaction
			try {
				prior_interaction = await this.AI.interactions.get(topic_id)
			} catch (err) {
				this.printDebug('INTERACTION_ERR', err, true)
				localStorage.removeItem(this.TOPIC_STORAGE_KEY)
				return {
					err: err,
				}
			}
			this.printDebug('PRIOR_INTERACTION', prior_interaction)
			return prior_interaction
		}
		return ''
	}

	printDebug(debug_name, debug_payload, is_err = false) {
		const debug_obj = {
			debug: debug_payload,
			name: debug_name,
			time: new Date().toISOString(),
		}
		if (this.DEBUG) {
			if (is_err) {
				console.error(debug_obj)
			} else {
				console.debug(debug_obj)
			}
		}
	}
}
