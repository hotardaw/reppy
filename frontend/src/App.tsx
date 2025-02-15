import React, { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  return (
    <div className="app-container" style={{
      backgroundColor: 'white',
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px'
    }}>
      <h1 style={{
        color: '#50514F',
        fontSize: '48px',
        marginBottom: '40px'
      }}>
        Reppy
      </h1>

      <button 
        style={{
          backgroundColor: '#F05D5E',
          color: 'white',
          padding: '15px 30px',
          border: 'none',
          borderRadius: '8px',
          fontSize: '18px',
          cursor: 'pointer',
          boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
          transition: 'transform 0.2s ease',
        }}
        onMouseOver={(e) => e.currentTarget.style.transform = 'scale(1.05)'}
        onMouseOut={(e) => e.currentTarget.style.transform = 'scale(1)'}
      >
        Sign in with Google
      </button>

      <p style={{
        color: '#1a73e8',
        fontSize: '14px',
        marginTop: '15px',
        cursor: 'pointer'
      }}>
        or log in/create an account
      </p>
    </div>
  );
}

export default App
