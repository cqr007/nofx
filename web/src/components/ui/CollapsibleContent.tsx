import React, { useState } from 'react'
import { Copy } from 'lucide-react'
import { t, type Language } from '../../i18n/translations'
import { copyWithToast } from '../../lib/clipboard'

interface CollapsibleContentProps {
  icon: React.ReactNode
  title: string
  titleColor: string
  content: string
  language: Language
  defaultExpanded?: boolean
}

export function CollapsibleContent({
  icon,
  title,
  titleColor,
  content,
  language,
  defaultExpanded = false,
}: CollapsibleContentProps) {
  const [expanded, setExpanded] = useState(defaultExpanded)

  return (
    <div className="mb-3">
      <button
        onClick={() => setExpanded(!expanded)}
        className="flex items-center gap-2 text-sm transition-colors"
        style={{ color: titleColor }}
      >
        <span className="font-semibold flex items-center gap-2">
          {icon} {title}
        </span>
        <span className="text-xs">
          {expanded ? t('collapse', language) : t('expand', language)}
        </span>
      </button>
      {expanded && (
        <div className="mt-2 relative">
          <button
            onClick={() => copyWithToast(content, t('copied', language))}
            className="absolute top-2 right-2 p-1.5 rounded hover:bg-gray-700 transition-colors z-10"
            title={t('copy', language)}
            style={{ background: 'rgba(43, 49, 57, 0.8)' }}
          >
            <Copy className="w-4 h-4" style={{ color: '#848E9C' }} />
          </button>
          <div
            className="rounded p-4 text-sm font-mono whitespace-pre-wrap max-h-96 overflow-y-auto"
            style={{
              background: '#0B0E11',
              border: '1px solid #2B3139',
              color: '#EAECEF',
            }}
          >
            {content}
          </div>
        </div>
      )}
    </div>
  )
}
