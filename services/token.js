const { createClient } = require('@supabase/supabase-js');

// Replace these with your values
const SUPABASE_URL = 'https://bwlnihqelepcrqmdkleh.supabase.co';
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImJ3bG5paHFlbGVwY3JxbWRrbGVoIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTQ0NzY1MjcsImV4cCI6MjA3MDA1MjUyN30.aGhnh6v68Lna5qXbeUzBBM2jiK_zQDkLpSPKhJQNPc0'
const EMAIL = 'aadityyadav1@gmil.com';
const PASSWORD = 'KKKK';

const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY);

async function loginAndGetToken() {
  const { data, error } = await supabase.auth.signInWithPassword({
    email: EMAIL,
    password: PASSWORD
  });

  if (error) {
    console.error('❌ Login failed:', error.message);
    return;
  }

  console.log('✅ Access Token:\n');
  console.log(data.session.access_token);
}

loginAndGetToken();
