import Gaia from '@gaia'
import { load } from '@dotenv'
import { parseArgs } from '@parse-args'

const LOCAL_STORAGE_DIR = '/tmp/localStorage'

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
	removeItem: (_key: string) => {},
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
		'nsfw',
		'print-env',
		'related',
		'transcript',
		'version',
	],
	string: [
		'ask',
		'env',
		'read',
	],
})

await load({
	envPath: FLAGS.env,
	export: true,
})

const API_KEY = Deno.env.get('GAIA_API_KEY')
const ATTACHMENTS: Attachment[] = []
const INTELLECT = Deno.env.get('GAIA_INTELLECT')
const MODEL = Deno.env.get('GAIA_MODEL')
const NSFW = Deno.env.get('GAIA_CENSOR_POLICY')
const PERSONA = Deno.env.get('GAIA_PERSONA')
const VERSION = 'v0.6.0 (main)'

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

if (NSFW) {
	AI.CENSOR_POLICY = NSFW
}

if (PERSONA) {
	AI.PERSONA = PERSONA
}

if (FLAGS.version) {
	AI.STDOUT(VERSION)
	Deno.exit()
}

if (FLAGS['print-env']) {
	AI.printEnv()
	Deno.exit()
}

if (FLAGS.echo) {
	AI.echo()
	Deno.exit()
}

if (FLAGS.transcript) {
	AI.printTranscript()
	Deno.exit()
}

if (FLAGS.nsfw) {
	AI.CENSOR_POLICY = 'OFF'
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
	user_prompt: String(FLAGS.ask),
}

const result = await AI.ask({
	...ask_obj,
	attachments: ask_obj.attachments as unknown as never[],
	transcript: undefined,
})

if (result.err) {
	Deno.exit(1)
}
