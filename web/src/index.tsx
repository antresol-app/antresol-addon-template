// The addon's web part: it registers its contributions when the app loads
// it (admins only in contract 1). React, the router and the SDK are the
// app's own (the Vite plugin maps the imports); build the UI from `ui`.
import { useEffect, useState } from 'react'
import { addonApi, registerAddon, ui, useAddonT } from '@antresol/addon-sdk'

const ID = 'template'
const api = addonApi(ID)

const strings = {
  ru: { title: 'Шаблон', hello: 'Привет, {name}', failed: 'Не удалось связаться с расширением' },
  en: { title: 'Template', hello: 'Hello, {name}', failed: 'Could not reach the addon' },
}
// labels and titles are fixed when the bundle registers: the app's language then
const lang = document.documentElement.lang === 'en' ? 'en' : 'ru'

function Page() {
  const t = useAddonT(ID)
  const [name, setName] = useState<string | null>(null)
  const [failed, setFailed] = useState(false)
  useEffect(() => {
    api.get<{ name: string }>('hello').then((r) => setName(r.name), () => setFailed(true))
  }, [])
  return (
    <section style={{ display: 'grid', gap: 'var(--s-4)' }}>
      {failed && <ui.Notice tone="destructive">{t('failed')}</ui.Notice>}
      {name !== null && <h1 style={{ font: 'var(--t-large-title)', margin: 0 }}>{t('hello', { name })}</h1>}
    </section>
  )
}

registerAddon(ID, {
  routes: [{ path: 'page', title: strings[lang].title, element: Page }],
  manageTiles: [{ label: strings[lang].title, icon: 'puzzle', to: 'page' }],
  strings,
})
