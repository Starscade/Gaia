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

	appendTranscript({
		data = '',
		is_agent = false,
		topic_id = '',
		type = 'text',
	} = {}) {
		const transcript = this.getTranscript()
		transcript.data.push({
			data: data,
			is_agent: is_agent,
			time: new Date().toISOString(),
			topic_id: topic_id,
			type: type,
		})
		this.printDebug('APPEND_TRANSCRIPT', transcript)
		this.setTranscript(transcript.data)
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
					{
						voice: 'Sulafat',
					},
				],
				thinking_level: this.INTELLECT,
			},
			model: this.MODELS[modality],
			response_modalities: [modality],
			stream: true,
			system_instruction: this.PERSONA,
		}

		this.appendTranscript({
			data: interaction_obj.input[0].text,
		})

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

				this.appendTranscript({
					data: template_attachment.data,
				})

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

		const transcript_buffer = {
			data: '',
			is_agent: true,
			type: 'text',
		}

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
							transcript_buffer.data += event.delta.text
							this.STDOUT({
								data: event.delta.text,
								type: event.delta.type,
							})
							break
						default:
							if (event.delta.data) {
								transcript_buffer.data += event.delta.data
								transcript_buffer.type = event.delta.type
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

		this.appendTranscript(transcript_buffer)

		return {
			err: null,
		}
	}

	forgetTranscript() {
		localStorage.removeItem(this.TOPIC_STORAGE_KEY)
		localStorage.removeItem(this.TRANSCRIPT_STORAGE_KEY)
		return null
	}

	async getEcho() {
		const transcript = await this.getTranscript()
		if (transcript.data && transcript.data.length > 0) {
			const final_entry = transcript.data[transcript.data.length - 1]
			return {
				data: final_entry.data,
				type: final_entry.type,
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

	getTranscript() {
		const return_obj = {
			data: [],
			err: null,
		}
		const raw_transcript = localStorage.getItem(this.TRANSCRIPT_STORAGE_KEY)
		if (raw_transcript) {
			try {
				const transcript = JSON.parse(raw_transcript)
				return_obj.data = transcript
			} catch (err) {
				return_obj.err = err
			}
		}
		this.printDebug('GET_TRANSCRIPT', return_obj)
		return return_obj
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

	setTranscript(transcript = []) {
		const err_obj = {
			err: null,
		}
		this.printDebug('SET_TRANSCRIPT', transcript)
		let json_transcript
		try {
			json_transcript = JSON.stringify(transcript)
		} catch (err) {
			err_obj.err = err
			return err_obj
		}
		localStorage.setItem(
			this.TRANSCRIPT_STORAGE_KEY,
			json_transcript,
		)
		return err_obj
	}
}
