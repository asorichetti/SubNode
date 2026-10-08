import React, { useEffect } from 'react'

const BrokenPage: React.FC = () => {
  useEffect(() => {
    // Simulate console errors
    console.error('This is a test error')
    console.warn('This is a test warning')
    
    // Simulate a failed network request
    fetch('/api/nonexistent-endpoint-that-will-fail')
      .catch(() => {})
    
    // Hide the page-content marker to break pagestate gate
    const marker = document.querySelector('[data-testid="page-content"]')
    if (marker) {
      (marker as HTMLElement).setAttribute('data-testid', 'broken-marker')
    }
    
    // Hide main element to break pagestate gate
    const main = document.querySelector('main')
    if (main) {
      (main as HTMLElement).style.display = 'none'
    }
  }, [])

  return (
    <div data-testid="broken-page" className="text-center py-24">
      <h1 className="text-3xl font-bold text-red-500 mb-4">Broken Page</h1>
      <p className="text-[#8b929a]">This page is intentionally broken to test UX gates.</p>
      <div className="mt-8 space-y-2">
        <div className="h-8 bg-red-900/20 rounded p-4">
          Missing content area
        </div>
      </div>
    </div>
  )
}

export default BrokenPage
