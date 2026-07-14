import Gaia from '@gaia'
import { load } from '@dotenv'
import { parseArgs } from '@parse-args'

const LOCAL_STORAGE_DIR = Deno.env.get('GAIA_STORAGE_DIR') ??
	'/tmp/localStorage'

try {
	Deno.mkdirSync(LOCAL_STORAGE_DIR)
} catch (err) {
	if (!(err instanceof Deno.errors.AlreadyExists)) {
		throw err
	}
}

const localStorageMock = {
	getItem: (key: string) => {
		let local_storage = ''
		try {
			local_storage = Deno.readTextFileSync(`${LOCAL_STORAGE_DIR}/${key}.json`)
		} catch (err) {
			if (err instanceof Deno.errors.NotFound) {
				Deno.writeTextFileSync(`${LOCAL_STORAGE_DIR}/${key}.json`, '')
			} else {
				throw err
			}
		}
		return local_storage
	},
	setItem: (key: string, value: string) => {
		return Deno.writeTextFileSync(`${LOCAL_STORAGE_DIR}/${key}.json`, value)
	},
	removeItem: (key: string) => {
		return Deno.remove(`${LOCAL_STORAGE_DIR}/${key}.json`)
	},
	clear: () => {},
	key: (_index: number) => null,
	length: 0,
}

Object.defineProperty(globalThis, 'localStorage', {
	value: localStorageMock,
	writable: true,
	configurable: true,
	enumerable: true,
})

interface Ask {
	attachments: Attachment[]
	modalities?: string[]
	preserve_context: boolean
	user_prompt: string
}

interface Attachment {
	data: string
	mime_type: string
}

const FLAGS = parseArgs(Deno.args, {
	boolean: [
		'draw',
		'echo',
		'forget',
		'environment',
		'related',
		'transcript',
		'version',
	],
	string: [
		'env',
		'read',
	],
})

const PROMPT_ARG = FLAGS._.length > 0 ? String(FLAGS._[FLAGS._.length - 1]) : ''
let PROMPT_PAD = ''
let PROMPT_STDIN = ''

if (!Deno.stdin.isTerminal()) {
	const raw = await new Response(Deno.stdin.readable).arrayBuffer()
	const STDIN = new TextDecoder().decode(raw).trim()

	if (STDIN.length > 0) {
		PROMPT_STDIN = STDIN
	}
}

if (PROMPT_ARG || PROMPT_STDIN) {
	PROMPT_PAD = '\n\n'
}

const USER_PROMPT = [
	PROMPT_ARG,
	PROMPT_STDIN,
].join(PROMPT_PAD).trim()

await load({
	envPath: FLAGS.env,
	export: true,
})

const API_KEY = Deno.env.get('GAIA_API_KEY')
const ATTACHMENTS: Attachment[] = []
const PERSONA = Deno.env.get('GAIA_PERSONA') ??
	'Your name is Gaia. You are a CLI tool. You respond exclusively in plaintext code snippets that can be executed (or compiled) as is. Never format your responses using markdown. If no language is specified, write code in POSIX-compliant sh (or PostgreSQL if dealing with SQL). Otherwise, write the code in the language that the user mentions. Always use the most portable shell syntax (e.g. the oldest, most widely supported). Always use the newest syntax if dealing with other languages. Never use node.js: use Deno instead. Prefer tab indentation to spaces. Never introduce yourself.'
const VERSION = 'v0.8.12 (dev)'

const AI = new Gaia({
	api_key: API_KEY,
	print_function: (text) => {
		Deno.stdout.writeSync(new TextEncoder().encode(text))
	},
})

if (PERSONA) {
	AI.PERSONA = PERSONA
}

if (FLAGS.version) {
	console.info(VERSION)
	Deno.exit()
}

if (FLAGS.environment) {
	const json_env = AI.getEnv()
	const current_environment = Object.entries(json_env).map(([key, value]) =>
		`${key.toUpperCase()}=${JSON.stringify(value)}`
	).join('\n')
	console.info(current_environment)
	Deno.exit()
}

if (FLAGS.forget) {
	AI.forgetTranscript()
	Deno.exit()
}

if (FLAGS.echo) {
	const final_thought = await AI.getEcho()
	console.info(final_thought)
	Deno.exit()
}

if (FLAGS.transcript) {
	const transcript = await AI.getTranscript()
	console.info(JSON.stringify(transcript))
	Deno.exit()
}

if (FLAGS.read) {
	const path = FLAGS.read
	try {
		const data = Deno.readFileSync(path)
		const b64 = btoa(String.fromCharCode(...data))
		const mime_type = path.endsWith('.png')
			? 'image/png'
			: path.endsWith('.jpg') || path.endsWith('.jpeg')
			? 'image/jpeg'
			: 'text/plain'
		ATTACHMENTS.push({
			mime_type,
			data: b64,
		})
	} catch (err) {
		console.error(`ERR: ${(err as Error).message}`)
		Deno.exit(1)
	}
}

if (!API_KEY) {
	console.error('ERR: GAIA_API_KEY is missing!')
	Deno.exit(1)
}

const ask_obj: Ask = {
	attachments: ATTACHMENTS,
	preserve_context: FLAGS.related,
	user_prompt: USER_PROMPT,
}

if (FLAGS.draw) {
	ask_obj.modalities = [
		'image',
	]
}

const result = await AI.ask({
	...ask_obj,
	attachments: ask_obj.attachments as unknown as never[],
})

if (result.err) {
	Deno.exit(1)
}

console.log()
