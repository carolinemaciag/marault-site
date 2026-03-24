// Marault Chatbot
const MARAULT_CONTEXT = 'Helper for Marault Intelligence data consulting';
let chatHistory = [];

const SERVICE_PAGES = {
  'data-visibility-audit': { url: '/services/data-visibility-audit', name: 'Data Visibility Audit' },
  'revenue-customer-analytics': { url: '/services/revenue-customer-analytics', name: 'Revenue & Customer Analytics' },
  'executive-dashboards': { url: '/services/executive-dashboards-reporting', name: 'Executive Dashboards' },
  'forecasting': { url: '/services/forecasting-decision-modeling', name: 'Forecasting & Decision Modeling' },
  'private-client': { url: '/services/private-client-analytics', name: 'Private Client Analytics' },
  'custom-build': { url: '/services/custom-website-build', name: 'Custom Website Build' },
  'template-build': { url: '/services/template-based-build', name: 'Template-Based Build' },
  'website-redesign': { url: '/services/website-redesign', name: 'Website Redesign' },
  'ux-ui-design': { url: '/services/ux-ui-design', name: 'UX/UI Design' },
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
    const handleToggleClick = function(e) {
      if (e) {
        e.preventDefault();
        e.stopPropagation();
      }
      container.classList.remove('chatbot-closed');
      container.classList.add('chat-open');
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

  if (closeBtn) {
    const handleClose = function(e) {
      e.preventDefault();
      e.stopPropagation();
      container.classList.add('chatbot-closed');
      container.classList.remove('chat-open');
      clearChatMemory();
    };
    closeBtn.addEventListener('click', handleClose);
    closeBtn.addEventListener('touchend', handleClose);
  }

  if (form) {
    form.addEventListener('submit', sendMessage);
  }
}


async function sendMessage(event) {
  if (event && event.preventDefault) {
    event.preventDefault();
  }
  const input = document.getElementById('chat-input');
  const message = input.value.trim();
  if (!message) return;

  addMessage(message, 'user');
  input.value = '';

  let response = checkGreeting(message) 
  || checkDataProblem(message)
  || checkPricing(message)
  || checkExperience(message)        // ✅ ADD
  || checkDifferentiation(message)   // ✅ ADD
  || checkGoodbye(message) 
  || checkQualifications(message) 
  || checkSecurity(message) 
  || checkTimeline(message) 
  || checkWhatDoDo(message) 
  || checkServiceDescription(message) 
  || checkServiceRecommendation(message) 
  || checkCompanyInfo(message) 
  || checkInappropriate(message) 
  || checkOffTopic(message);

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
    console.error('Chat fetch error:', e);
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

function checkDataProblem(msg) {
  const m = msg.toLowerCase();

  if (
    /messy data|data.*messy|clean.*data|fix.*data|data.*fix|data.*problem|bad data|data quality|data is wrong|data is broken/.test(m)
  ) {
    return '<p>Yes — that’s exactly what we specialize in.</p>' +
           '<p>We help clean, structure, and transform messy data into a clear, decision-ready system.</p>' +
           '<p>This typically starts with our ' +
           createServiceLink('data-visibility-audit', 'Data Visibility Audit') +
           '.</p>';
  }

  return null;
}

function checkPricing(msg) {
  const m = msg.toLowerCase();

  if (
    m.includes('price') ||
    m.includes('cost') ||
    m.includes('how much') ||
    m.includes('pricing') ||
    m.includes('fee') ||
    m.includes('budget') ||
    m.includes('rate')
  ) {
    return '<p>Our pricing structure depends on the scale and scope of your project.</p>' +
           '<p>We begin with a conversation to understand your data, goals, and where we can create the most value.</p>' +
           '<p>' + createServiceLink('inquire', 'Get in touch with us for a consultation and tailored quote') + '.</p>';
  }

  return null;
}

function checkExperience(msg) {
  const m = msg.toLowerCase();

  if (
    m.includes('worked with') ||
    m.includes('clients') ||
    m.includes('companies') ||
    m.includes('experience') ||
    m.includes('who have you worked with') ||
    m.includes('types of clients') ||
    m.includes('large companies')
  ) {
    return '<p>We work with both growing companies and more established organizations, across a range of industries.</p>' +
           '<p>Our focus is less on company size, and more on solving meaningful data problems — whether that’s improving visibility, building decision systems, or creating more clarity for leadership.</p>' +
           '<p>' + createServiceLink('inquire', 'We’re happy to discuss your specific situation here') + '.</p>';
  }

  return null;
}

function checkDifferentiation(msg) {
  const m = msg.toLowerCase();

  if (
    m.includes('why choose') ||
    m.includes('why you') ||
    m.includes('difference') ||
    m.includes('better than') ||
    m.includes('vs') ||
    m.includes('versus') ||
    m.includes('competitor') ||
    m.includes('compare')
  ) {
    return '<p>We focus on clarity and decision-making — not just building models or dashboards for the sake of it.</p>' +
           '<p>Our work is designed to be practical, interpretable, and directly tied to how you run your business.</p>' +
           '<p>We also take a more tailored approach than many firms — every engagement is built around your specific data, goals, and constraints.</p>' +
           '<p>In short, we don’t just deliver analysis — we help you make better decisions with confidence.</p>' +
           '<p>' + createServiceLink('approach', 'You can learn more about how we work here') + '.</p>';
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

  const auditMatch = m === 'audit' || m === '1' || /data visibility audit|data audit|what.*audit|tell.*audit|describe.*audit|info.*audit/.test(m);
  const revenueMatch = m === 'revenue' || m === '2' || /revenue|customer analytics|what.*revenue|tell.*revenue|describe.*revenue|info.*revenue/.test(m);
  const executiveMatch = m === 'executive' || m === '3' || /executive|dashboards|reporting|what.*executive|tell.*executive|describe.*executive|describe.*dashboards|info.*executive/.test(m);
  const forecastingMatch = /forecasting|decision modeling|what.*forecast|tell.*forecast|describe.*forecast|info.*forecast/.test(m);
  const privateMatch = /private client|what.*private|tell.*private|describe.*private|info.*private/.test(m);
  const customMatch = /custom website|custom build|what.*custom|tell.*custom|describe.*custom|info.*custom/.test(m);
  const templateMatch = /template|template.*build|what.*template|tell.*template|describe.*template|info.*template/.test(m);
  const uxMatch = m === 'ux' || m === 'ui' || /ux|ui|ux\/ui|ux ui|ux design|what.*ux|tell.*ux|describe.*ux|describe.*design|info.*ux|what.*design|tell.*design|info.*design/.test(m);

  if (auditMatch) return '<p><strong>Data Visibility Audit</strong> - We assess your current data infrastructure, identify gaps, and provide a roadmap for better data access and transparency. ' + createServiceLink('data-visibility-audit', 'Learn more') + '</p>';
  if (revenueMatch) return '<p><strong>Revenue & Customer Analytics</strong> - Unlock deep insights into customer behavior patterns and revenue drivers. ' + createServiceLink('revenue-customer-analytics', 'Learn more') + '</p>';
  if (executiveMatch) return '<p><strong>Executive Dashboards & Reporting</strong> - Real-time dashboards and custom reports give leadership clear visibility into the KPIs that matter most. ' + createServiceLink('executive-dashboards', 'Learn more') + '</p>';
  if (forecastingMatch) return '<p><strong>Forecasting & Decision Modeling</strong> - Our predictive models forecast future trends and test scenarios to inform strategic planning. ' + createServiceLink('forecasting', 'Learn more') + '</p>';
  if (privateMatch) return '<p><strong>Private Client Analytics</strong> - Specialized analytics for high-net-worth individuals and private businesses. ' + createServiceLink('private-client', 'Learn more') + '</p>';
  if (customMatch) return '<p><strong>Custom Website Builds</strong> - Fully custom-built websites tailored to your exact specifications and business goals. ' + createServiceLink('custom-build', 'Learn more') + '</p>';
  if (templateMatch) return '<p><strong>Template-Based Builds</strong> - Fast, cost-effective website solutions built on proven templates. ' + createServiceLink('template-build', 'Learn more') + '</p>';
  if (uxMatch) return '<p><strong>UX/UI Design</strong> - We design intuitive, beautiful user experiences that engage customers and drive conversions. ' + createServiceLink('ux-ui-design', 'Learn more') + '</p>';

  return null;
}

function checkInappropriate(msg) {
  const m = msg.toLowerCase();
  const inappropriatePatterns = /\b(racist|sexist|porn|xxx|rape|harass|lewd|perverted|nsfw)\b/i;
  if (inappropriatePatterns.test(m)) {
    return '<p>I don\'t feel comfortable discussing that topic. I\'m here to help with Marault Intelligence services. ' + createServiceLink('services', 'Let\'s talk about what we do') + ' instead!</p>';
  }
  return null;
}

function checkOffTopic(msg) {
  const m = msg.toLowerCase();
  const offTopicPatterns = /what.*eat|lunch|dinner|breakfast|pizza|burger|coffee|weather|sports|movie|funny|joke|cat|dog|pet|music|song|game|hobby|vacation|travel|recipe|cooking/;
  if (offTopicPatterns.test(m) && !(/marault|data|service|business|analytics|audit|dashboard|build|design|website/i.test(m))) {
    return '<p>That\'s a fun question, but I\'m specifically here to help with Marault Intelligence services.</p><p>' + createServiceLink('services', 'Explore our services') + '</p>';
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