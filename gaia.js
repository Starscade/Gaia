export default class {
	constructor({
		api_key = '',
		censor_policy = 'BLOCK_ONLY_HIGH',
		intellect = 'MINIMAL',
		model = 'gemini-flash-lite-latest',
		persona =
			`Your name is Gaia. You are a CLI tool. You respond exclusively in plaintext code snippets that can be executed (or compiled) as is. Never format your responses using markdown. If no language is specified, write code in POSIX-compliant sh (or PostgreSQL if dealing with SQL). Otherwise, write the code in the language that the user mentions. Always use the most portable shell syntax (e.g. the oldest, most widely supported). Always use the newest syntax if dealing with other languages. Never use node.js: use Deno instead. Prefer tab indentation to spaces. Never introduce yourself.`,
		print_function = (text) => {
			console.log(text)
		},
	} = {}) {
		this.API_KEY = api_key
		this.CENSOR_POLICY = censor_policy
		this.INTELLECT = intellect
		this.MODEL = model
		this.PERSONA = persona
		this.STDOUT = print_function
		this.TRANSCRIPT_STORAGE_KEY = 'GAIA_TRANSCRIPT'
	}

	async ask({
		attachments = [],
		preserve_context = false,
		transcript = [],
		user_prompt = '',
	} = {}) {
		if (!user_prompt) {
			if (attachments.length === 0) {
				const err_msg = 'No prompt provided!'
				console.error('ERR:', err_msg)
				return {
					data: null,
					err: err_msg,
				}
			}
		}

		if (preserve_context) {
			const saved = localStorage.getItem(this.TRANSCRIPT_STORAGE_KEY)
			if (saved) {
				transcript = JSON.parse(saved)
			}
		}

		// console.info('TRANSCRIPT:', transcript)

		const user_parts = [{ text: user_prompt }]
		for (const file of attachments) {
			user_parts.push({
				inlineData: {
					mimeType: file.mime_type,
					data: file.data,
				},
			})
		}

		const contents = [
			...transcript.map((msg) => ({
				role: msg.role === 'model' ? 'model' : 'user',
				parts: msg.parts || [{ text: msg.text }],
			})),
		]

		// console.info('CONTENTS:', contents)

		contents.push({ role: 'user', parts: user_parts })
		localStorage.setItem(this.TRANSCRIPT_STORAGE_KEY, JSON.stringify(contents))

		const json_body = {
			contents: contents,
			generationConfig: {
				thinkingConfig: { thinkingLevel: this.INTELLECT },
			},
			safetySettings: [
				{
					category: 'HARM_CATEGORY_DANGEROUS_CONTENT',
					threshold: this.CENSOR_POLICY,
				},
				{
					category: 'HARM_CATEGORY_HARASSMENT',
					threshold: this.CENSOR_POLICY,
				},
				{
					category: 'HARM_CATEGORY_HATE_SPEECH',
					threshold: this.CENSOR_POLICY,
				},
				{
					category: 'HARM_CATEGORY_SEXUALLY_EXPLICIT',
					threshold: this.CENSOR_POLICY,
				},
			],
			system_instruction: { parts: [{ text: this.PERSONA }] },
			tools: [{
				google_search: {},
			}],
		}

		const response = await fetch(
			`https://generativelanguage.googleapis.com/v1beta/models/${this.MODEL}:streamGenerateContent`,
			{
				body: JSON.stringify(json_body),
				headers: {
					'Content-Type': 'application/json',
					'x-goog-api-key': this.API_KEY,
				},
				method: 'POST',
			},
		)

		if (!response.ok) {
			const json_response = await response.json()
			console.error('ERR:', json_response)
			return {
				data: null,
				err: json_response,
			}
		}

		const reader = response.body.getReader()
		const decoder = new TextDecoder()
		let buffer = ''
		let lastIndex = 0
		let fullText = ''

		while (true) {
			const { done, value } = await reader.read()
			if (done) break

			buffer += decoder.decode(value, { stream: true })

			const attempt = buffer.trim().endsWith(']') ? buffer : buffer + ']'

			try {
				const data = JSON.parse(attempt)
				for (let i = lastIndex; i < data.length; i++) {
					const text = data[i].candidates?.[0].content?.parts?.[0].text
					if (text) {
						fullText += text
						this.STDOUT(text)
					}
				}
				lastIndex = data.length
			} catch (_err) {
				// Silently ignore partial chunks until stream is complete.
			}
		}

		contents.push({ role: 'model', parts: [{ text: fullText }] })
		localStorage.setItem(this.TRANSCRIPT_STORAGE_KEY, JSON.stringify(contents))

		return {
			data: contents,
			err: null,
		}
	}

	echo() {
		const raw_transcript = localStorage.getItem(this.TRANSCRIPT_STORAGE_KEY)
		const transcript = JSON.parse(raw_transcript)
		const final_thought = transcript[transcript.length - 1]?.parts[0]?.text
		this.STDOUT(final_thought)
		return {
			data: final_thought,
			err: null,
		}
	}

	printEnv() {
		const json_env = {
			GAIA_API_KEY: this.API_KEY,
			GAIA_CENSOR_POLICY: this.CENSOR_POLICY,
			GAIA_INTELLECT: this.INTELLECT,
			GAIA_MODEL: this.MODEL,
			GAIA_PERSONA: this.PERSONA,
		}
		const current_environment = Object.entries(json_env).map(([key, value]) =>
			`${key.toUpperCase()}=${JSON.stringify(value)}`
		).join('\n')
		this.STDOUT(current_environment)
		return {
			data: json_env,
			err: null,
		}
	}

	printTranscript() {
		const transcript = localStorage.getItem(this.TRANSCRIPT_STORAGE_KEY)
		this.STDOUT(transcript)
		return {
			data: transcript,
			err: null,
		}
	}
}
