const API_URL = 'http://localhost:8080/v1/auth';

class AuthService {
  static async register(userData) {
    const response = await fetch(`${API_URL}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(userData)
    });
    if (!response.ok) throw new Error('Registration failed');
    return await response.json();
  }

  static async login(credentials) {
    const response = await fetch(`${API_URL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials)
    });
    if (!response.ok) throw new Error('Login failed');
    return await response.json();
  }

  static async loginPhone(phone) {
    const response = await fetch(`${API_URL}/login-phone`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ phone })
    });
    if (!response.ok) throw new Error('Phone login failed');
    return await response.json();
  }

  static async verifyOTP(phone, code) {
    const response = await fetch(`${API_URL}/verify-otp`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ phone, code })
    });
    if (!response.ok) throw new Error('OTP verification failed');
    return await response.json();
  }

  static async checkAuth(token) {
    const response = await fetch(`${API_URL}/validate`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      }
    });
    if (!response.ok) throw new Error('Token validation failed');
    return await response.json();
  }
}

export default AuthService;
