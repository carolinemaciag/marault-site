// Marault Chatbot
const MARAULT_CONTEXT = 'Helper for Marault Intelligence data consulting';
let chatHistory = [];

const SERVICE_PAGES = {
  'data-visibility-audit': { url: '/data-visibility-audit', name: 'Data Visibility Audit' },
  'revenue-customer-analytics': { url: '/revenue', name: 'Revenue & Customer Analytics' },
  'executive-dashboards': { url: '/executive-dashboards-reporting', name: 'Executive Dashboards' },
  'forecasting': { url: '/forecasting-decision-modeling', name: 'Forecasting & Decision Modeling' },
  'private-client': { url: '/private-client-analytics', name: 'Private Client Analytics' },
  'custom-build': { url: '/custom-website-build', name: 'Custom Website Build' },
  'template-build': { url: '/template-based-build', name: 'Template-Based Build' },
  'website-redesign': { url: '/website-redesign', name: 'Website Redesign' },
  'ux-ui-design': { url: '/ux-ui-design', name: 'UX/UI Design' },
  'services': { url: '/services', name: 'Services' },
  'inquire': { url: '/inquire', name: 'Inquire' },
  'contact': { url: '/contact', name: 'Contact' },
  'approach': { url: '/approach', name: 'Approach' },
  'team': { url: '/executive-team', name: 'Team' },
  'philosophy': { url: '/philosophy', name: 'Philosophy' },
};

function createServiceLink(key, text) {
  const service = SERVICE_PAGES[key];
  if (!service) return text || key;
  return '<a href="' + service.url + '" target="_blank" class="chat-service-link">' + (text || service.name) + '</a>';
}

document.addEventListener('DOMContentLoaded', setupChatbot);

function setupChatbot() {
  const toggleBtn = document.getElementById('chat-toggle');
  const closeBtn = document.getElementById('close-chat');
  const container = document.getElementById('chatbot-container');
  const form = document.getElementById('chat-form');

  if (toggleBtn) {
    // Direct click handler with multiple event types for better Safari support
    const handleToggleClick = function(e) {
      if (e) {
        e.preventDefault();
        e.stopPropagation();
      }
      container.classList.remove('chatbot-closed');
      container.classList.add('chat-open');
      const saved = localStorage.getItem('marault_chat_history');
      if (saved) {
        chatHistory = JSON.parse(saved);
        renderChatHistory();
      }
      setTimeout(() => {
        const input = document.getElementById('chat-input');
        if (input) input.focus();
      }, 100);
    };
    
    toggleBtn.addEventListener('click', handleToggleClick);
    toggleBtn.addEventListener('touchend', handleToggleClick);
    toggleBtn.addEventListener('touchstart', (e) => {
      e.preventDefault();
    });
  }

  container.addEventListener('click', (e) => {
    const isMobile = window.innerWidth <= 480;
    const isNearBottom = e.clientY > window.innerHeight - 80;
    const isClosed = container.classList.contains('chatbot-closed');
    
    if (isMobile && isNearBottom && isClosed) {
      e.preventDefault();
      e.stopPropagation();
      container.classList.remove('chatbot-closed');
      container.classList.add('chat-open');
      const saved = localStorage.getItem('marault_chat_history');
      if (saved) {
        chatHistory = JSON.parse(saved);
        renderChatHistory();
      }
      setTimeout(() => {
        const input = document.getElementById('chat-input');
        if (input) input.focus();
      }, 100);
    }
  });

  // Also handle touchend for better mobile support
  container.addEventListener('touchend', (e) => {
    const isMobile = window.innerWidth <= 480;
    const touch = e.changedTouches[0];
    const isNearBottom = touch.clientY > window.innerHeight - 80;
    const isClosed = container.classList.contains('chatbot-closed');
    
    if (isMobile && isNearBottom && isClosed) {
      e.preventDefault();
      e.stopPropagation();
      container.classList.remove('chatbot-closed');
      container.classList.add('chat-open');
      const saved = localStorage.getItem('marault_chat_history');
      if (saved) {
        chatHistory = JSON.parse(saved);
        renderChatHistory();
      }
      setTimeout(() => {
        const input = document.getElementById('chat-input');
        if (input) input.focus();
      }, 100);
    }
  });

  if (closeBtn) {
    closeBtn.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      container.classList.add('chatbot-closed');
      container.classList.remove('chat-open');
      const isMobile = window.innerWidth <= 480;
      if (isMobile) {
        clearChatMemory();
      }
    });
    closeBtn.addEventListener('touchend', (e) => {
      e.preventDefault();
      e.stopPropagation();
      container.classList.add('chatbot-closed');
      container.classList.remove('chat-open');
      const isMobile = window.innerWidth <= 480;
      if (isMobile) {
        clearChatMemory();
      }
    });
  }

  if (form) {
    form.addEventListener('submit', sendMessage);
  }

  const saved = localStorage.getItem('marault_chat_history');
  if (saved) {
    chatHistory = JSON.parse(saved);
    renderChatHistory();
  }
}

async function sendMessage(event) {
  event.preventDefault();
  const input = document.getElementById('chat-input');
  const message = input.value.trim();
  if (!message) return;

  addMessage(message, 'user');
  input.value = '';

  let response = checkGreeting(message) || checkGoodbye(message) || checkQualifications(message) || checkSecurity(message) || checkTimeline(message) || checkWhatDoDo(message) || checkServiceDescription(message) || checkServiceRecommendation(message) || checkCompanyInfo(message) || checkInappropriate(message) || checkOffTopic(message);
  
  if (response) {
    setTimeout(() => addMessage(response, 'bot', true), 600);
    if (checkGoodbye(message)) setTimeout(() => { document.getElementById('chatbot-container').classList.add('chatbot-closed'); clearChatMemory(); }, 2100);
    return;
  }

  showTypingIndicator();
  try {
    const res = await fetch('/api/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message, history: chatHistory, context: MARAULT_CONTEXT }),
    });
    const data = await res.json();
    removeTypingIndicator();
    addMessage(data.reply, 'bot');
  } catch (e) {
    removeTypingIndicator();
    addMessage('Sorry, having trouble responding. Please try again or contact us.', 'bot');
  }
}

function checkGreeting(msg) {
  const m = msg.toLowerCase().trim();
  if (/^(hello|hi|hey|good\s+(morning|afternoon|evening)|what'?s\s+up|sup|yo|greetings)/.test(m)) {
    return ['Hey there! I\'m muh·ROH. I\'ll help you explore Marault\'s services!', 'Hello! I\'m muh·ROH. How can I assist you today?', 'Hi! I\'m muh·ROH. What brings you here?'][Math.floor(Math.random() * 3)];
  }
  return null;
}

function checkGoodbye(msg) {
  const m = msg.toLowerCase().trim();
  if (/^(bye|goodbye|see\s+you|take\s+care|thanks|that's\s+all|all\s+set)/.test(m)) {
    return ['Take care! Reach out anytime.', 'Thanks! Have a great day!', 'See you soon!'][Math.floor(Math.random() * 3)];
  }
  return null;
}

function checkQualifications(msg) {
  const m = msg.toLowerCase();
  if (/who are you|team|founder|caroline|lindsey|qualifications|who runs|runs this company|who's running|leadership/.test(m)) {
    return '<p><strong>Marault Intelligence</strong> is founded and run by two data science experts:</p><p><strong>Caroline Maciag</strong> - Deep Learning & Time Series specialist with MS from Northwestern</p><p><strong>Lindsey Chenault</strong> - AI & Data Integrity specialist with MS from Northwestern</p><p>Together, they bring elite training and expertise in data science and business intelligence. ' + createServiceLink('team', 'Meet the Team') + '</p>';
  }
  return null;
}

function checkSecurity(msg) {
  const m = msg.toLowerCase();
  if (/security|data security|privacy|compliance|gdpr|soc 2|encrypted|secure|protection|confidential/.test(m)) {
    return '<p><strong>Data Security & Compliance</strong></p><p>Data security is paramount to us. We follow industry best practices including encrypted data transmission, secure storage protocols, and strict confidentiality agreements. We\'re compliant with major regulations including GDPR and SOC 2.</p><p>For detailed security documentation and certifications, please ' + createServiceLink('contact', 'reach out via our contact page') + '.</p>';
  }
  return null;
}

function checkTimeline(msg) {
  const m = msg.toLowerCase();
  if (/how long|timeline|weeks|duration|how many weeks|how long does|timeframe|how quickly|rush|expedited/.test(m)) {
    return '<p><strong>Project Timeline</strong></p><p>Every project is unique based on scope and complexity. ' + createServiceLink('inquire', 'Let us discuss your timeline') + '.</p>';
  }
  return null;
}

function checkServiceRecommendation(msg) {
  const m = msg.toLowerCase();
  if (/which service|what service|help.*choose|not.*sure.*which/.test(m)) {
    return '<p>Here are our main services:</p><ul><li>' + createServiceLink('data-visibility-audit', 'Data Visibility Audit') + '</li><li>' + createServiceLink('executive-dashboards', 'Executive Dashboards') + '</li><li>' + createServiceLink('revenue-customer-analytics', 'Revenue & Customer Analytics') + '</li><li>' + createServiceLink('forecasting', 'Forecasting & Decision Modeling') + '</li></ul><p>Tell me more about your needs!</p>';
  }
  return null;
}

function checkCompanyInfo(msg) {
  const m = msg.toLowerCase();
  if (/^what is.*company|^about|^who is marault/.test(m)) {
    return '<p><strong>Marault Intelligence</strong> is a data consulting firm transforming data into actionable insights.</p><p><strong>Philosophy:</strong> Clarity in Analysis. Confidence in Action.</p><p>' + createServiceLink('services', 'Services') + ' | ' + createServiceLink('approach', 'Our Approach') + ' | ' + createServiceLink('team', 'Team') + '</p>';
  }
  return null;
}

function checkWhatDoDo(msg) {
  const m = msg.toLowerCase();
  if (/^what do you do|^what.*do.*you.*do|^services|^what do we offer|i need.*service|need service|what.*service|tell.*service|describe.*service/.test(m)) {
    return '<p>We offer 8 core services:</p><ol><li>' + createServiceLink('data-visibility-audit', 'Data Visibility Audit') + '</li><li>' + createServiceLink('revenue-customer-analytics', 'Revenue & Customer Analytics') + '</li><li>' + createServiceLink('executive-dashboards', 'Executive Dashboards & Reporting') + '</li><li>' + createServiceLink('forecasting', 'Forecasting & Decision Modeling') + '</li><li>' + createServiceLink('private-client', 'Private Client Analytics') + '</li><li>' + createServiceLink('custom-build', 'Custom Website Builds') + '</li><li>' + createServiceLink('template-build', 'Template-Based Builds') + '</li><li>' + createServiceLink('ux-ui-design', 'UX/UI Design') + '</li></ol><p>Tell me which service interests you!</p>';
  }
  return null;
}

function checkServiceDescription(msg) {
  const m = msg.toLowerCase().trim();
  
  // Check for various question patterns about services
  const auditMatch = m === 'audit' || m === '1' || /data visibility audit|data audit|what.*audit|tell.*audit|describe.*audit|info.*audit/.test(m);
  const revenueMatch = m === 'revenue' || m === '2' || /revenue|customer analytics|what.*revenue|tell.*revenue|describe.*revenue|info.*revenue/.test(m);
  const executiveMatch = m === 'executive' || m === '3' || /executive|dashboards|reporting|what.*executive|tell.*executive|describe.*executive|describe.*dashboards|info.*executive/.test(m);
  const forecastingMatch = /forecasting|decision modeling|what.*forecast|tell.*forecast|describe.*forecast|info.*forecast/.test(m);
  const privateMatch = /private client|what.*private|tell.*private|describe.*private|info.*private/.test(m);
  const customMatch = /custom website|custom build|what.*custom|tell.*custom|describe.*custom|info.*custom/.test(m);
  const templateMatch = /template|template.*build|what.*template|tell.*template|describe.*template|info.*template/.test(m);
  const uxMatch = m === 'ux' || m === 'ui' || /ux|ui|ux\/ui|ux ui|ux design|what.*ux|tell.*ux|describe.*ux|describe.*design|info.*ux|what.*design|tell.*design|info.*design/.test(m);
  
  if (auditMatch) {
    return '<p><strong>Data Visibility Audit</strong> - We assess your current data infrastructure, identify gaps, and provide a roadmap for better data access and transparency. This is the foundation for better decision-making. ' + createServiceLink('data-visibility-audit', 'Learn more') + '</p>';
  }
  if (revenueMatch) {
    return '<p><strong>Revenue & Customer Analytics</strong> - Unlock deep insights into customer behavior patterns and revenue drivers. We help you optimize pricing, identify high-value segments, and maximize profitability. ' + createServiceLink('revenue-customer-analytics', 'Learn more') + '</p>';
  }
  if (executiveMatch) {
    return '<p><strong>Executive Dashboards & Reporting</strong> - Real-time dashboards and custom reports give leadership clear visibility into the KPIs that matter most to your business. Make faster, data-driven decisions. ' + createServiceLink('executive-dashboards', 'Learn more') + '</p>';
  }
  if (forecastingMatch) {
    return '<p><strong>Forecasting & Decision Modeling</strong> - Our predictive models forecast future trends and test scenarios to inform strategic planning. Plan with confidence based on data, not guesswork. ' + createServiceLink('forecasting', 'Learn more') + '</p>';
  }
  if (privateMatch) {
    return '<p><strong>Private Client Analytics</strong> - Specialized analytics for high-net-worth individuals and private businesses. We create tailored solutions for complex financial situations and confidential analysis. ' + createServiceLink('private-client', 'Learn more') + '</p>';
  }
  if (customMatch) {
    return '<p><strong>Custom Website Builds</strong> - Fully custom-built websites tailored to your exact specifications and business goals. We combine beautiful design with powerful functionality and data integration. ' + createServiceLink('custom-build', 'Learn more') + '</p>';
  }
  if (templateMatch) {
    return '<p><strong>Template-Based Builds</strong> - Fast, cost-effective website solutions built on proven templates. Perfect for getting online quickly while maintaining professional design and core functionality. ' + createServiceLink('template-build', 'Learn more') + '</p>';
  }
  if (uxMatch) {
    return '<p><strong>UX/UI Design</strong> - We design intuitive, beautiful user experiences that engage customers and drive conversions. From wireframes to pixel-perfect designs, we create interfaces people love to use. ' + createServiceLink('ux-ui-design', 'Learn more') + '</p>';
  }
  
  return null;
}

function checkInappropriate(msg) {
  const m = msg.toLowerCase();
  // Check for hateful, perverted, or otherwise inappropriate content
  const inappropriatePatterns = /\b(hate|racist|sexist|slur|lewd|perverted|nsfw|adult|porn|xxx|sexual|harass|abuse|violence|kill|rape|assault|discriminat|offensive|degrad|retard|stupid|dumb|idiot|waste)\b/i;
  
  if (inappropriatePatterns.test(m)) {
    return '<p>I don\'t feel comfortable discussing that topic. However, I\'m here to help with anything related to Marault Intelligence and our business services. ' + createServiceLink('services', 'Let\'s talk about what we do') + ' instead!</p>';
  }
  return null;
}

function checkOffTopic(msg) {
  const m = msg.toLowerCase();
  // Check if it's clearly off-topic (not about services, company, or business)
  const offTopicPatterns = /what.*eat|lunch|dinner|breakfast|pizza|burger|coffee|weather|sports|movie|funny|joke|cat|dog|pet|music|song|game|hobby|vacation|travel|recipe|cooking/;
  
  if (offTopicPatterns.test(m) && !(/marault|data|service|business|analytics|audit|dashboard|build|design|website/i.test(m))) {
    return '<p>That\'s a fun question, but I\'m specifically here to help with Marault Intelligence services.</p><p>' + createServiceLink('services', 'Explore our services') + ' or let me know how we can help with your data needs.</p>';
  }
  return null;
}

function addMessage(msg, sender, isHTML) {
  const container = document.getElementById('chat-messages');
  const div = document.createElement('div');
  div.className = 'chat-message ' + sender + '-message';
  const content = document.createElement('div');
  content.className = 'message-content';
  if (isHTML) content.innerHTML = msg; else content.textContent = msg;
  div.appendChild(content);
  container.appendChild(div);
  chatHistory.push({ sender, message: msg, isHTML });
  localStorage.setItem('marault_chat_history', JSON.stringify(chatHistory));
  container.scrollTop = container.scrollHeight;
}

function showTypingIndicator() {
  const container = document.getElementById('chat-messages');
  const div = document.createElement('div');
  div.id = 'typing-indicator';
  div.className = 'chat-message bot-message';
  const content = document.createElement('div');
  content.className = 'message-content';
  content.innerHTML = '<span class="typing"></span><span class="typing"></span><span class="typing"></span>';
  div.appendChild(content);
  container.appendChild(div);
  container.scrollTop = container.scrollHeight;
}

function removeTypingIndicator() {
  const el = document.getElementById('typing-indicator');
  if (el) el.remove();
}

function renderChatHistory() {
  const container = document.getElementById('chat-messages');
  if (chatHistory.length > 0) {
    container.innerHTML = '';
    chatHistory.forEach(item => addMessage(item.message, item.sender, item.isHTML));
  }
}

function clearChatMemory() {
  chatHistory = [];
  localStorage.removeItem('marault_chat_history');
  const container = document.getElementById('chat-messages');
  container.innerHTML = '';
  const div = document.createElement('div');
  div.className = 'chat-message bot-message';
  const content = document.createElement('div');
  content.className = 'message-content';
  content.textContent = 'Hi! I\'m muh·ROH, your Marault Intelligence assistant. Welcome! I\'m here to help you explore our services and answer any questions. How can I assist you today?';
  div.appendChild(content);
  container.appendChild(div);
}
