import { Octokit } from "octokit";

const GITHUB_TOKEN = Deno.env.get("GITHUB_TOKEN");
const PORT = parseInt(Deno.env.get("PORT") || "8000");

if (!GITHUB_TOKEN) {
  console.error("❌ GITHUB_TOKEN environment variable is required");
  Deno.exit(1);
}

const octokit = new Octokit({ auth: GITHUB_TOKEN });

interface UserStats {
  username: string;
  name: string;
  totalStars: number;
  totalForks: number;
  totalRepos: number;
  publicRepos: number;
  privateRepos: number;
  followers: number;
  following: number;
  totalCommits: number;
  totalIssues: number;
  totalPRs: number;
}

interface LanguageStats {
  [language: string]: number;
}

// Cache to avoid hitting rate limits
const cache = new Map<string, { data: any; timestamp: number }>();
const CACHE_DURATION = 3600000; // 1 hour in milliseconds

function getFromCache(key: string) {
  const cached = cache.get(key);
  if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
    return cached.data;
  }
  return null;
}

function setCache(key: string, data: any) {
  cache.set(key, { data, timestamp: Date.now() });
}

async function getUserStats(username: string, includePrivate = false): Promise<UserStats> {
  const cacheKey = `user-stats-${username}-${includePrivate ? 'private' : 'public'}`;
  const cached = getFromCache(cacheKey);
  if (cached) return cached;

  // Check if this is the authenticated user
  let isAuthenticatedUser = false;
  const cachedAuthUser = getFromCache("auth-user");

  if (cachedAuthUser) {
    isAuthenticatedUser = cachedAuthUser.login.toLowerCase() === username.toLowerCase();
  } else {
    try {
      const { data: authUser } = await octokit.rest.users.getAuthenticated();
      isAuthenticatedUser = authUser.login.toLowerCase() === username.toLowerCase();
      setCache("auth-user", authUser);
    } catch (error) {
      // Not authenticated or token doesn't have user scope
      console.warn("Could not get authenticated user:", error);
    }
  }

  // Get user info
  const { data: user } = await octokit.rest.users.getByUsername({ username });

  // Get repositories - use different endpoints based on whether we want private repos
  let repos;
  if (includePrivate && isAuthenticatedUser) {
    // Get all repos including private ones for authenticated user
    repos = await octokit.paginate(octokit.rest.repos.listForAuthenticatedUser, {
      per_page: 100,
      affiliation: 'owner',
      visibility: 'all', // Include both public and private
    });
  } else {
    // Get only public repos
    repos = await octokit.paginate(octokit.rest.repos.listForUser, {
      username,
      per_page: 100,
    });
  }

  // Calculate stats
  const totalStars = repos.reduce((sum, repo) => sum + (repo.stargazers_count ?? 0), 0);
  const totalForks = repos.reduce((sum, repo) => sum + (repo.forks_count ?? 0), 0);
  const publicRepos = repos.filter(repo => !repo.private).length;
  const privateRepos = repos.filter(repo => repo.private).length;

  // Get commit count (approximate from recent activity)
  let totalCommits = 0;
  try {
    const commits = await octokit.paginate("GET /search/commits", {
      q: `author:${username}`,
      per_page: 100,
    });
    totalCommits = commits.length;
  } catch (error) {
    console.warn("Could not fetch commits:", error);
  }

  // Get issues and PRs
  let totalIssues = 0;
  let totalPRs = 0;
  try {
    const issues = await octokit.paginate("GET /search/issues", {
      q: `author:${username} type:issue`,
      per_page: 100,
    });
    totalIssues = issues.length;

    const prs = await octokit.paginate("GET /search/issues", {
      q: `author:${username} type:pr`,
      per_page: 100,
    });
    totalPRs = prs.length;
  } catch (error) {
    console.warn("Could not fetch issues/PRs:", error);
  }

  const stats: UserStats = {
    username: user.login,
    name: user.name || user.login,
    totalStars,
    totalForks,
    totalRepos: repos.length,
    publicRepos: includePrivate ? publicRepos : user.public_repos,
    privateRepos: includePrivate ? privateRepos : 0,
    followers: user.followers,
    following: user.following,
    totalCommits,
    totalIssues,
    totalPRs,
  };

  setCache(cacheKey, stats);
  return stats;
}

async function getLanguageStats(username: string, includePrivate = false): Promise<LanguageStats> {
  const cacheKey = `lang-stats-${username}-${includePrivate ? 'private' : 'public'}`;
  const cached = getFromCache(cacheKey);
  if (cached) return cached;

  // Check if this is the authenticated user
  let isAuthenticatedUser = false;
  const cachedAuthUser = getFromCache("auth-user");

  if (cachedAuthUser) {
    isAuthenticatedUser = cachedAuthUser.login.toLowerCase() === username.toLowerCase();
  } else {
    try {
      const { data: authUser } = await octokit.rest.users.getAuthenticated();
      isAuthenticatedUser = authUser.login.toLowerCase() === username.toLowerCase();
      setCache("auth-user", authUser);
    } catch (error) {
      // Not authenticated or token doesn't have user scope
      console.warn("Could not get authenticated user:", error);
    }
  }

  // Get repositories - use different endpoints based on whether we want private repos
  let repos;
  if (includePrivate && isAuthenticatedUser) {
    repos = await octokit.paginate(octokit.rest.repos.listForAuthenticatedUser, {
      per_page: 100,
      affiliation: 'owner',
      visibility: 'all',
    });
  } else {
    repos = await octokit.paginate(octokit.rest.repos.listForUser, {
      username,
      per_page: 100,
    });
  }

  const languages: LanguageStats = {};

  for (const repo of repos) {
    if (repo.fork) continue; // Skip forked repos

    try {
      const { data: repoLanguages } = await octokit.rest.repos.listLanguages({
        owner: repo.owner.login,
        repo: repo.name,
      });

      for (const [lang, bytes] of Object.entries(repoLanguages)) {
        languages[lang] = (languages[lang] || 0) + bytes;
      }
    } catch (error) {
      console.warn(`Could not fetch languages for ${repo.name}:`, error);
    }
  }

  // Convert bytes to percentages
  const total = Object.values(languages).reduce((sum, bytes) => sum + bytes, 0);
  const percentages: LanguageStats = {};
  for (const [lang, bytes] of Object.entries(languages)) {
    percentages[lang] = parseFloat(((bytes / total) * 100).toFixed(2));
  }

  // Sort by percentage
  const sorted = Object.fromEntries(
    Object.entries(percentages).sort(([, a], [, b]) => b - a)
  );

  setCache(cacheKey, sorted);
  return sorted;
}

function generateStatsSVG(stats: UserStats): string {
  const width = 465;
  const height = 195;
  
  return `
<svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .header { font: 600 18px 'Segoe UI', Ubuntu, Sans-Serif; fill: #fff; }
      .stat { font: 400 14px 'Segoe UI', Ubuntu, Sans-Serif; fill: #9f9f9f; }
      .stat-value { font: 600 14px 'Segoe UI', Ubuntu, Sans-Serif; fill: #fff; }
      .icon { fill: #79ff97; }
    </style>
  </defs>
  
  <rect width="${width}" height="${height}" fill="#151515" rx="4.5"/>
  
  <text x="25" y="35" class="header">${stats.name}'s GitHub Stats</text>
  
  <!-- Stats -->
  <g transform="translate(0, 55)">
    <!-- Left column -->
    <g transform="translate(25, 0)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Stars: <tspan class="stat-value">${stats.totalStars}</tspan></text>
    </g>
    
    <g transform="translate(25, 30)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M5 3.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm0 2.122a2.25 2.25 0 10-1.5 0v.878A2.25 2.25 0 005.75 8.5h1.5v2.128a2.251 2.251 0 101.5 0V8.5h1.5a2.25 2.25 0 002.25-2.25v-.878a2.25 2.25 0 10-1.5 0v.878a.75.75 0 01-.75.75h-4.5A.75.75 0 015 6.25v-.878z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Forks: <tspan class="stat-value">${stats.totalForks}</tspan></text>
    </g>
    
    <g transform="translate(25, 60)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M2 2.5A2.5 2.5 0 014.5 0h8.75a.75.75 0 01.75.75v12.5a.75.75 0 01-.75.75h-2.5a.75.75 0 110-1.5h1.75v-2h-8a1 1 0 00-.714 1.7.75.75 0 01-1.072 1.05A2.495 2.495 0 012 11.5v-9zm10.5-1V9h-8c-.356 0-.694.074-1 .208V2.5a1 1 0 011-1h8z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Repos: <tspan class="stat-value">${stats.totalRepos}</tspan></text>
    </g>
    
    <g transform="translate(25, 90)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Commits: <tspan class="stat-value">${stats.totalCommits}</tspan></text>
    </g>
    
    <!-- Right column -->
    <g transform="translate(260, 0)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
      </svg>
      <text x="25" y="12" class="stat">Total PRs: <tspan class="stat-value">${stats.totalPRs}</tspan></text>
    </g>
    
    <g transform="translate(260, 30)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 9.5a1.5 1.5 0 100-3 1.5 1.5 0 000 3z"/>
        <path d="M8 0a8 8 0 100 16A8 8 0 008 0zM1.5 8a6.5 6.5 0 1113 0 6.5 6.5 0 01-13 0z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Issues: <tspan class="stat-value">${stats.totalIssues}</tspan></text>
    </g>
    
    <g transform="translate(260, 60)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M5.5 3.5a2 2 0 100 4 2 2 0 000-4zM2 5.5a3.5 3.5 0 115.898 2.549 5.507 5.507 0 013.034 4.084.75.75 0 11-1.482.235 4.001 4.001 0 00-7.9 0 .75.75 0 01-1.482-.236A5.507 5.507 0 013.102 8.05 3.49 3.49 0 012 5.5zM11 4a.75.75 0 100 1.5 1.5 1.5 0 01.666 2.844.75.75 0 00-.416.672v.352a.75.75 0 00.574.73c1.2.289 2.162 1.2 2.522 2.372a.75.75 0 101.434-.44 5.01 5.01 0 00-2.56-3.012A3 3 0 0011 4z"/>
      </svg>
      <text x="25" y="12" class="stat">Followers: <tspan class="stat-value">${stats.followers}</tspan></text>
    </g>
    
    <g transform="translate(260, 90)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M5.5 3.5a2 2 0 100 4 2 2 0 000-4zM2 5.5a3.5 3.5 0 115.898 2.549 5.507 5.507 0 013.034 4.084.75.75 0 11-1.482.235 4.001 4.001 0 00-7.9 0 .75.75 0 01-1.482-.236A5.507 5.507 0 013.102 8.05 3.49 3.49 0 012 5.5zM11 4a.75.75 0 100 1.5 1.5 1.5 0 01.666 2.844.75.75 0 00-.416.672v.352a.75.75 0 00.574.73c1.2.289 2.162 1.2 2.522 2.372a.75.75 0 101.434-.44 5.01 5.01 0 00-2.56-3.012A3 3 0 0011 4z"/>
      </svg>
      <text x="25" y="12" class="stat">Following: <tspan class="stat-value">${stats.following}</tspan></text>
    </g>
  </g>
</svg>`.trim();
}

function generateLanguagesSVG(languages: LanguageStats): string {
  const width = 300;
  const entries = Object.entries(languages).slice(0, 8); // Top 8 languages
  const height = 45 + entries.length * 40;
  
  const colors: { [key: string]: string } = {
    JavaScript: "#f1e05a",
    TypeScript: "#3178c6",
    Python: "#3572A5",
    Java: "#b07219",
    C: "#555555",
    "C++": "#f34b7d",
    "C#": "#178600",
    Go: "#00ADD8",
    Rust: "#dea584",
    Ruby: "#701516",
    PHP: "#4F5D95",
    Swift: "#ffac45",
    Kotlin: "#A97BFF",
    Dart: "#00B4AB",
    HTML: "#e34c26",
    CSS: "#563d7c",
    Shell: "#89e051",
  };
  
  return `
<svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .header { font: 600 18px 'Segoe UI', Ubuntu, Sans-Serif; fill: #fff; }
      .lang-name { font: 400 14px 'Segoe UI', Ubuntu, Sans-Serif; fill: #9f9f9f; }
      .percentage { font: 400 12px 'Segoe UI', Ubuntu, Sans-Serif; fill: #9f9f9f; }
    </style>
  </defs>
  
  <rect width="${width}" height="${height}" fill="#151515" rx="4.5"/>
  
  <text x="25" y="35" class="header">Most Used Languages</text>
  
  <g transform="translate(0, 50)">
    ${entries.map(([lang, percent], i) => `
      <g transform="translate(25, ${i * 40})">
        <circle cx="6" cy="8" r="6" fill="${colors[lang] || "#858585"}"/>
        <text x="20" y="12" class="lang-name">${lang}</text>
        <text x="${width - 80}" y="12" class="percentage">${percent}%</text>
      </g>
    `).join("")}
  </g>
</svg>`.trim();
}

async function handler(req: Request): Promise<Response> {
  const url = new URL(req.url);
  const username = url.searchParams.get("username") || "elyor04";
  const includePrivate = url.searchParams.get("private") === "true";
  
  // CORS headers
  const headers = new Headers({
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Methods": "GET, OPTIONS",
    "Access-Control-Allow-Headers": "Content-Type",
  });

  if (req.method === "OPTIONS") {
    return new Response(null, { headers });
  }

  try {
    if (url.pathname === "/stats") {
      const stats = await getUserStats(username, includePrivate);
      const svg = generateStatsSVG(stats);
      
      headers.set("Content-Type", "image/svg+xml");
      headers.set("Cache-Control", "public, max-age=3600");
      
      return new Response(svg, { headers });
    } 
    else if (url.pathname === "/languages") {
      const languages = await getLanguageStats(username, includePrivate);
      const svg = generateLanguagesSVG(languages);
      
      headers.set("Content-Type", "image/svg+xml");
      headers.set("Cache-Control", "public, max-age=3600");
      
      return new Response(svg, { headers });
    }
    else if (url.pathname === "/api/stats") {
      const stats = await getUserStats(username, includePrivate);
      
      headers.set("Content-Type", "application/json");
      headers.set("Cache-Control", "public, max-age=3600");
      
      return new Response(JSON.stringify(stats, null, 2), { headers });
    }
    else if (url.pathname === "/api/languages") {
      const languages = await getLanguageStats(username, includePrivate);
      
      headers.set("Content-Type", "application/json");
      headers.set("Cache-Control", "public, max-age=3600");
      
      return new Response(JSON.stringify(languages, null, 2), { headers });
    }
    else {
      headers.set("Content-Type", "application/json");
      return new Response(
        JSON.stringify({ error: "Not found" }),
        { status: 404, headers }
      );
    }
  } catch (error: any) {
    console.error("Error:", error);
    headers.set("Content-Type", "application/json");
    return new Response(
      JSON.stringify({ error: error.message }),
      { status: 500, headers }
    );
  }
}

if (import.meta.main) {
  Deno.cron("update elyor04 stats", "*/5 * * * *", async () => {
    await getUserStats("elyor04", true);
    await getLanguageStats("elyor04", true);
  })

  console.log(`🚀 Server running on http://localhost:${PORT}`);
  Deno.serve({ port: PORT }, handler);
}
