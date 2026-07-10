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
	preserve_context: boolean
	transcript?: string
	user_prompt: string
}

interface Attachment {
	mime_type: string
	data: string
}

const FLAGS = parseArgs(Deno.args, {
	boolean: [
		'echo',
		'forget',
		'nsfw',
		'print-env',
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
const INTELLECT = Deno.env.get('GAIA_INTELLECT')
const MODEL = Deno.env.get('GAIA_MODEL')
const NSFW = Deno.env.get('GAIA_NSFW')
const PERSONA = Deno.env.get('GAIA_PERSONA') ??
	'Your name is Gaia. You are a CLI tool. You respond exclusively in plaintext code snippets that can be executed (or compiled) as is. Never format your responses using markdown. If no language is specified, write code in POSIX-compliant sh (or PostgreSQL if dealing with SQL). Otherwise, write the code in the language that the user mentions. Always use the most portable shell syntax (e.g. the oldest, most widely supported). Always use the newest syntax if dealing with other languages. Never use node.js: use Deno instead. Prefer tab indentation to spaces. Never introduce yourself.'
const VERSION = 'v0.7.40 (main)'

const AI = new Gaia({
	api_key: API_KEY,
	print_function: (text) => {
		Deno.stdout.writeSync(new TextEncoder().encode(text))
	},
})

if (INTELLECT) {
	AI.INTELLECT = INTELLECT
}

if (MODEL) {
	AI.MODEL = MODEL
}

AI.NSFW = FLAGS.nsfw || [
	'true',
	'1',
	'yes',
	'on',
].includes(
	NSFW?.toLowerCase().trim() ?? '',
)

if (PERSONA) {
	AI.PERSONA = PERSONA
}

if (FLAGS.version) {
	AI.STDOUT(`${VERSION}\n`)
	Deno.exit()
}

if (FLAGS['print-env']) {
	AI.printEnv()
	AI.STDOUT('\n')
	Deno.exit()
}

if (FLAGS.forget) {
	const result = AI.forgetTranscript()
	AI.STDOUT(result)
	Deno.exit()
}

if (FLAGS.echo) {
	AI.echo()
	Deno.exit()
}

if (FLAGS.transcript) {
	AI.printTranscript()
	AI.STDOUT('\n')
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

const result = await AI.ask({
	...ask_obj,
	attachments: ask_obj.attachments as unknown as never[],
	transcript: undefined,
})

if (result.err) {
	Deno.exit(1)
}

AI.STDOUT('\n')
