export interface LumberJackConfig {
  baseUrl: string;
  token?: string;
  timeout?: number;
  retryAttempts?: number;
}

export interface RequestConfig extends RequestInit {
  timeout?: number;
  retryAttempts?: number;
}

export interface AuthConfig {
  token: string;
  expiresAt: number;
}

export interface LoginResponse {
  token: string;
}

export interface TimeTrackingEntry {
  start_time: string;
  end_time: string;
  duration: number;
}

export interface EventSummary {
  event_id: string;
  status: 'pending' | 'ongoing' | 'finished';
  start_time: string;
  end_time: string;
  entries_count: number;
  metadata: Record<string, any>;
}

export type Permission = 'read' | 'write' | 'admin';

export type NodeType = 'branch' | 'leaf';

export type EventStatus = 'pending' | 'ongoing' | 'finished';

export interface User {
  id: string;
  username: string;
  email: string;
  organization?: string;
  phone?: string;
}

export interface Node {
  id: string;
  name: string;
  type: NodeType;
  path: string;
  children?: Node[];
  parents?: string[];
  metadata?: Record<string, any>;
}

export interface Event {
  event_id: string;
  path: string;
  status: EventStatus;
  start_time: string;
  end_time?: string;
  entries: EventEntry[];
  metadata?: Record<string, any>;
}

export interface EventEntry {
  timestamp: string;
  content: string;
  user_id: string;
  metadata?: Record<string, any>;
}

export interface UserAssignment {
  path: string;
  assignee_id: string;
  permission: Permission;
}

export interface ServerSettings {
  organization: string;
  server_port: string;
  dashboard_url: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  organization?: string;
  phone?: string;
}

export interface APIResponse<T> {
  data?: T;
  error?: string;
}

export class LumberJack {
  private baseUrl: string;
  private token: string | null;

  constructor(config: LumberJackConfig) {
    this.baseUrl = config.baseUrl.replace(/\/$/, '');
    this.token = config.token || null;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers: {
        ...headers,
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'API request failed');
    }

    return response.json();
  }

  async login(username: string, password: string): Promise<void> {
    const response = await this.request<LoginResponse>('/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
    this.token = response.token;
  }

  async startTimeTracking(path: string): Promise<void> {
    await this.request('/time/start', {
      method: 'POST',
      body: JSON.stringify({ path }),
    });
  }

  async stopTimeTracking(path: string): Promise<void> {
    await this.request('/time/stop', {
      method: 'POST',
      body: JSON.stringify({ path }),
    });
  }

  async getTimeTracking(path: string): Promise<TimeTrackingEntry[]> {
    return this.request('/time', {
      method: 'GET',
      body: JSON.stringify({ path }),
    });
  }

  async startEvent(path: string, eventId: string, metadata?: Record<string, any>): Promise<void> {
    await this.request('/events/start', {
      method: 'POST',
      body: JSON.stringify({ path, event_id: eventId, metadata }),
    });
  }

  async planEvent(path: string, eventId: string, startTime: Date, endTime: Date, metadata?: Record<string, any>): Promise<void> {
    await this.request('/events/plan', {
      method: 'POST',
      body: JSON.stringify({
        path,
        event_id: eventId,
        start_time: startTime.toISOString(),
        end_time: endTime.toISOString(),
        metadata,
      }),
    });
  }

  async appendToEvent(path: string, eventId: string, content: string, metadata?: Record<string, any>): Promise<void> {
    await this.request('/events/append', {
      method: 'POST',
      body: JSON.stringify({
        path,
        event_id: eventId,
        content,
        metadata,
      }),
    });
  }

  async endEvent(path: string, eventId: string): Promise<void> {
    await this.request('/events/end', {
      method: 'POST',
      body: JSON.stringify({ path, event_id: eventId }),
    });
  }

  async getForest(): Promise<any> {
    return this.request('/forest', {
      method: 'GET',
    });
  }

  async getTree(): Promise<any> {
    return this.request('/forest/tree', {
      method: 'GET',
    });
  }
}
