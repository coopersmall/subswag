import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'
import * as tsToZod from 'ts-to-zod'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const projectDir = path.resolve(__dirname, '..')

function generateValidation(): string {
  const typesDir = path.join(projectDir, 'src', 'types')

  if (!fs.existsSync(typesDir)) {
    fs.mkdirSync(typesDir, { recursive: true })
  }

  function collectTypeScriptFiles(dir: string): Array<{
    inputPath: string
    outputPath: string
    name: string
  }> {
    const files: Array<{
      inputPath: string
      outputPath: string
      name: string
    }> = []
    const items = fs.readdirSync(dir, { withFileTypes: true })

    for (const item of items) {
      if (
        item.isFile() &&
        item.name.endsWith('.generated.ts') &&
        !item.name.endsWith('.zod.ts')
      ) {
        const name = item.name.replace('.generated.ts', '')
        const inputPath = path.join(dir, item.name)
        const outputPath = path.join(dir, `${name}.generated.zod.ts`)

        files.push({
          inputPath,
          outputPath,
          name,
        })
      }
    }

    return files
  }

  function processFiles(
    files: Array<{ inputPath: string; outputPath: string; name: string }>
  ): void {
    // Create input/output mappings for cross-file references
    const inputOutputMappings = files.map((file) => ({
      input: `./${file.name}.generated.ts`,
      output: `./${file.name}.generated.zod.ts`,
    }))

    const customJSDocFormatTypes = {
      id: 'number',
    }

    for (const file of files) {
      try {
        const sourceText = fs.readFileSync(file.inputPath, 'utf8')

        const result = tsToZod.generate({
          sourceText,
          keepComments: false,
          inputOutputMappings,
          customJSDocFormatTypes,
        })

        if (result.errors.length > 0) {
          console.warn(`Warnings for ${file.inputPath}:`, result.errors)
        }

        const schemaContent = result.getZodSchemasFile(file.inputPath)

        fs.writeFileSync(file.outputPath, schemaContent)
        console.log(`Generated Zod schema for ${file.name}`)

        if (result.hasCircularDependencies) {
          console.warn(
            `Note: ${file.name} has circular dependencies that were resolved`
          )
        }
      } catch (error) {
        console.error(`Error processing ${file.name}:`, error)
      }
    }
  }

  const typeScriptFiles = collectTypeScriptFiles(typesDir)
  processFiles(typeScriptFiles)
  return `Validation schemas generated successfully.`
}

generateValidation()
