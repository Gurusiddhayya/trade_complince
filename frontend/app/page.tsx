'use client'

import { useEffect, useMemo, useState } from 'react'

const API = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1'

type Product = { code: string; name: string; category: string }
type Tx = { id: string; reference_no: string; product_type: string; product_code: string; transaction_nature?: 'Goods'|'Services / Software'; status: string; currency: string; declared_value: number; counterparty: string; country: string; description: string; contract_value: number; invoice_value: number; realized_value: number; outstanding_value: number; difference: number; remarks?: string }
type PaymentRecord = { id:string; transactionId:string; direction:'Inward'|'Outward'; reference:string; amount:number; date:string; currency:string; invoiceRef:string; remarks?:string }
type AuditEvent = { id:string; at:string; action:string; area:string; reference:string; detail:string }

type Screening = { id: string; entity: string; role: string; country: string; date: string; result: 'Clear' | 'Review' | 'Alert'; detail: string; source: string }

const PRODUCTS: Product[] = [
  ['EXP_GOODS','Export of Goods','Export'],['EXP_SERVICE','Software / Service Export','Export'],['IMP_GOODS','Import of Goods / Services','Import'],['HSS','High Sea Sale','Import'],['ILC','Import LC','Trade Finance'],['BG','Bank Guarantee','Trade Finance'],['EPC','EPC / PCFC','Export Finance'],['EBC','Export Bill Collection','Export'],['TC','Buyer / Supplier Credit','Trade Credit'],['MTT','Merchanting Trade Transaction','Trade'],['FDI','FDI','Investment'],['ODI','ODI','Investment'],['ECB','ECB','Borrowing']
].map(([code,name,category])=>({code,name,category}))

const SCREENING: Screening[] = [
  {id:'SCR-00041',entity:'Global Tech Ltd.',role:'Buyer',country:'USA',date:'19 Sep 2026 10:24',result:'Clear',detail:'No potential match found',source:'Sanctions + Watchlists'},
  {id:'SCR-00040',entity:'Sunrise Trading',role:'Supplier',country:'China',date:'18 Sep 2026 16:32',result:'Clear',detail:'No potential match found',source:'Sanctions + Watchlists'},
  {id:'SCR-00039',entity:'Metro Exports',role:'Buyer',country:'Germany',date:'17 Sep 2026 13:15',result:'Review',detail:'Potential low-confidence name similarity',source:'Watchlist'},
  {id:'SCR-00038',entity:'Prime Construction',role:'Supplier',country:'Singapore',date:'16 Sep 2026 11:03',result:'Clear',detail:'No potential match found',source:'Sanctions + Watchlists'},
  {id:'SCR-00037',entity:'ABC Holdings',role:'Counterparty',country:'UAE',date:'15 Sep 2026 09:47',result:'Clear',detail:'No potential match found',source:'Sanctions + Watchlists'}
]

const MOCK_TX: Tx[] = [
  {id:'1',reference_no:'EXP-2026-000125',product_type:'Export of Goods',product_code:'EXP',status:'With Bank',currency:'USD',declared_value:125000,counterparty:'Global Tech Ltd.',country:'USA',description:'Export of engineering components',contract_value:125000,invoice_value:125000,realized_value:60000,outstanding_value:65000,difference:0},
  {id:'2',reference_no:'IMP-2026-000098',product_type:'Import of Goods',product_code:'IMP',status:'In Progress',currency:'USD',declared_value:85000,counterparty:'Sunrise Trading',country:'China',description:'Import of industrial equipment',contract_value:85000,invoice_value:85000,realized_value:85000,outstanding_value:0,difference:0},
  {id:'3',reference_no:'HSS-2026-000076',product_type:'High Sea Sale',product_code:'HSS',status:'Attention Required',currency:'USD',declared_value:250000,counterparty:'ABC Holdings',country:'UAE',description:'High sea sale record',contract_value:250000,invoice_value:250000,realized_value:0,outstanding_value:250000,difference:0},
  {id:'4',reference_no:'ILC-2026-000065',product_type:'Import LC',product_code:'ILC',status:'With Bank',currency:'USD',declared_value:560000,counterparty:'Metro Exports',country:'Germany',description:'Import LC against purchase order',contract_value:560000,invoice_value:560000,realized_value:0,outstanding_value:560000,difference:0},
  {id:'5',reference_no:'BG-2026-000054',product_type:'Bank Guarantee',product_code:'BG',status:'Completed',currency:'USD',declared_value:100000,counterparty:'Prime Construction',country:'Singapore',description:'Performance guarantee',contract_value:100000,invoice_value:0,realized_value:0,outstanding_value:0,difference:0}
]

const money = (n:number,c='USD') => new Intl.NumberFormat('en-IN',{style:'currency',currency:c==='INR'?'INR':'USD',maximumFractionDigits:0}).format(n||0)

export default function Home(){
  const [tab,setTab]=useState('Dashboard')
  const [products,setProducts]=useState<Product[]>(PRODUCTS)
  const [tx,setTx]=useState<Tx[]>(MOCK_TX)
  const [dash,setDash]=useState<any>(null)
  const [selected,setSelected]=useState<Tx|null>(null)
  const [filter,setFilter]=useState('All')
  const [wizard,setWizard]=useState<Product|null>(null)
  const [screening,setScreening]=useState<Screening[]>(SCREENING)
  const [screeningFilter,setScreeningFilter]=useState('All')
  const [search,setSearch]=useState('')
  const [notice,setNotice]=useState('')
  const [importNature,setImportNature]=useState<'Goods'|'Services / Software'>('Goods')
  const [paymentRecords,setPaymentRecords]=useState<PaymentRecord[]>([])
  const [auditLog,setAuditLog]=useState<AuditEvent[]>([])
  const [editingTx,setEditingTx]=useState<Tx|null>(null)
  const [securityLocked,setSecurityLocked]=useState(false)

  const recordAudit=(action:string,area:string,reference:string,detail:string)=>setAuditLog(p=>[{id:`AUD-${Date.now()}`,at:new Date().toLocaleString(),action,area,reference,detail},...p])
  const updateTransaction=(updated:Tx)=>{setTx(prev=>prev.map(t=>t.id===updated.id?updated:t));setSelected(updated);setEditingTx(null);recordAudit('Updated','Transaction',updated.reference_no,'Customer edited transaction details or remarks.');setNotice('Transaction updated. The change was recorded in the local activity history.')}


  useEffect(()=>{
    Promise.all([
      fetch(API+'/products').then(r=>r.ok?r.json():null).catch(()=>null),
      fetch(API+'/transactions').then(r=>r.ok?r.json():null).catch(()=>null),
      fetch(API+'/dashboard').then(r=>r.ok?r.json():null).catch(()=>null)
    ]).then(([p,t,d])=>{
      if(Array.isArray(p)&&p.length) setProducts(p)
      if(Array.isArray(t)&&t.length) setTx(t)
      if(d) setDash(d)
    })
  },[])

  const metrics = {
    transactions: dash?.transactions ?? 42,
    actions: dash?.attention ?? 8,
    requirements: dash?.requirements?.open ?? 12,
    payments: 5,
    screening: screening.filter(x=>x.result!=='Clear').length
  }

  const visible = useMemo(()=>{
    let list = filter==='All' ? tx : tx.filter(t=>t.product_type===filter)
    if(search.trim()){
      const q=search.toLowerCase()
      list=list.filter(t=>[t.reference_no,t.product_type,t.counterparty,t.country].join(' ').toLowerCase().includes(q))
    }
    return list
  },[tx,filter,search])

  const filteredScreening = screeningFilter==='All' ? screening : screening.filter(x=>x.result===screeningFilter)

  const start = (p:Product)=>{ setWizard(p); setNotice('') }
  const createDemo = ()=>{
    const p=wizard||products[0]
    const n=tx.length+126
    const nature=p.code==='IMP_GOODS'?importNature:undefined
    const item:Tx={id:String(Date.now()),reference_no:`${p.code}-${new Date().getFullYear()}-${String(n).padStart(6,'0')}`,product_type:p.name,product_code:p.code,transaction_nature:nature,status:'Draft',currency:'USD',declared_value:0,counterparty:'',country:'',description:nature==='Services / Software'?'Import of services/software':'',contract_value:0,invoice_value:0,realized_value:0,outstanding_value:0,difference:0}
    setTx([item,...tx]); setWizard(null); setTab('Transactions'); setSelected(item); setNotice('Demo transaction created. Next: add transaction details, documents, requirements and payments.')
  }

  const runScreening=(entity:string)=>{
    const item:Screening={id:`SCR-${String(Date.now()).slice(-5)}`,entity,role:'Counterparty',country:'—',date:'Just now',result:'Clear',detail:'No potential match found in this demo screening run',source:'Demo screening source'}
    setScreening([item,...screening]); setNotice(`Screening completed for ${entity}. Result: Clear — no potential match found in this demo.`)
  }

  return <main>
    <header className="topbar">
      <div className="brand"><div className="brandMark">✓</div><div><h1>Trade Compliance Cockpit</h1><p>Your Trade Records. Our Support.</p></div></div>
      <div className="globalSearch"><span>⌕</span><input value={search} onChange={e=>setSearch(e.target.value)} placeholder="Search by transaction, invoice, buyer or reference…"/></div>
      <div className="user"><span className="bell">◉</span><div className="avatar">JD</div><div><b>Demo User</b><small>Sample Company · Demo data</small></div></div>
    </header>

    <div className="shell">
      <aside className="sidebar">
        {['Dashboard','Transactions','Documents','Payments & Reconciliation','Declarations','Bank Requirements','Buyer / Supplier Requirements','Action Center','Screening','Transaction 360','Reports','Help Me'].map(x=><button key={x} className={tab===x?'sideActive':''} onClick={()=>setTab(x)}>{x==='Dashboard'?'⌂':x==='Transactions'?'↔':x==='Documents'?'▤':x==='Payments & Reconciliation'?'₹':x==='Declarations'?'▤':x==='Bank Requirements'?'▣':x==='Buyer / Supplier Requirements'?'♙':x==='Action Center'?'!':x==='Screening'?'◈':x==='Transaction 360'?'◎':x==='Reports'?'▥':'?' }<span>{x}</span>{x==='Action Center'&&<em>3</em>}</button>)}
        <div className="sidebarBottom"><button onClick={()=>setTab('Settings')}>⚙<span>Security & Settings</span></button><small>Trade Compliance Cockpit<br/>V2.23 · Customer-friendly + security hardening</small></div>
      </aside>

      <section className="content">
        {notice && <div className="toast" onClick={()=>setNotice('')}>{notice}<b>×</b></div>}
        {tab==='Dashboard' && <Dashboard metrics={metrics} tx={tx} screening={screening} onOpen={(t)=>{setSelected(t);setTab('Transaction 360')}} onNew={()=>start(products[0])} onScreen={()=>setTab('Screening')} onUpload={()=>setNotice('Document upload flow is ready for the next implementation pass.')} onPayment={()=>setNotice('Payment recording flow is ready for the next implementation pass.')} />}
        {tab==='Transactions' && <Transactions tx={visible} products={products} filter={filter} setFilter={setFilter} onOpen={(t)=>{setSelected(t);setTab('Transaction 360')}} />}
        {tab==='Screening' && <ScreeningPage rows={filteredScreening} filter={screeningFilter} setFilter={setScreeningFilter} run={runScreening} />}
        {tab==='Declarations' && <DeclarationsPage tx={tx} onNotice={setNotice} />}
        {tab==='Transaction 360' && <Transaction360 tx={selected||tx[0]} onBack={()=>setTab('Transactions')} onEdit={()=>setEditingTx(selected||tx[0])} onRemark={(remark)=>{const t=selected||tx[0]; if(!t)return; updateTransaction({...t,remarks:remark})}} />}
        {tab==='Payments & Reconciliation' && <PaymentsWorkspace tx={tx} records={paymentRecords} onAdd={(r)=>{setPaymentRecords(prev=>[r,...prev]);setTx(prev=>prev.map(t=>{if(t.id!==r.transactionId||r.direction!=='Inward')return t;const realized=t.realized_value+r.amount;const base=t.invoice_value||t.declared_value;return {...t,realized_value:realized,outstanding_value:Math.max(0,base-realized),difference:Math.max(0,realized-base)}}));setNotice('Payment recorded in this browser demo. It has not been sent to a bank.')}} />}
        {tab==='Settings' && <SecuritySettings locked={securityLocked} onLock={()=>setSecurityLocked(true)} auditLog={auditLog} />}
        {['Documents','Bank Requirements','Buyer / Supplier Requirements','Action Center','Reports','Help Me'].includes(tab) && <Placeholder tab={tab} tx={tx} onOpen={(t)=>{setSelected(t);setTab('Transaction 360')}} />}
      </section>
    </div>

    {wizard && <div className="modal"><div className="modalCard"><div className="modalHead"><div><h2>Start a new transaction</h2><p>Choose the product first. We will ask only what this transaction needs.</p></div><button onClick={()=>setWizard(null)}>×</button></div>{wizard.code==='IMP_GOODS' && <label style={{display:'block',marginBottom:14}}>Import transaction type<select value={importNature} onChange={e=>setImportNature(e.target.value as 'Goods'|'Services / Software')}><option>Goods</option><option>Services / Software</option></select><small style={{display:'block',color:'#718096',marginTop:5}}>Goods may require Bill of Entry details. Services/software use service invoice and payment records; no BOE field is assumed.</small></label>}<div className="productSelect">{products.map(p=><button className={wizard.code===p.code?'selected':''} key={p.code} onClick={()=>setWizard(p)}><b>{p.name}</b><small>{p.category}</small></button>)}</div><div className="modalActions"><button onClick={()=>setWizard(null)}>Cancel</button><button className="primary" onClick={createDemo}>Create transaction</button></div></div></div>}
    {editingTx && <EditTransactionModal tx={editingTx} onCancel={()=>setEditingTx(null)} onSave={updateTransaction} />}
    {securityLocked && <div className="modal"><div className="modalCard lockCard"><h2>Workspace locked</h2><p>This prototype lock protects the visible workspace only. Production access must use server-side authentication and session controls.</p><button className="primary" onClick={()=>setSecurityLocked(false)}>Unlock demo workspace</button></div></div>}
    {selected && tab==='Transaction 360' && <div className="srOnly">{selected.reference_no}</div>}
  </main>
}

function Dashboard({metrics,tx,screening,onOpen,onNew,onScreen,onUpload,onPayment}:{metrics:any,tx:Tx[],screening:Screening[],onOpen:(t:Tx)=>void,onNew:()=>void,onScreen:()=>void,onUpload:()=>void,onPayment:()=>void}){
  const total=tx.reduce((a,t)=>a+t.declared_value,0)
  return <>
    <div className="pageHead"><div><h2>Welcome to your trade workspace</h2><p>Here’s what’s happening with your trade transactions today.</p></div><div className="date">Sample workspace · Demo data</div></div>
    <div className="metricGrid">
      <Metric icon="▣" label="Total Transactions" value={metrics.transactions} note="↑ 12% vs. last 7 days" />
      <Metric icon="!" label="Open Actions" value={metrics.actions} note="3 high priority · 5 medium" />
      <Metric icon="▤" label="Bank Requirements" value={metrics.requirements} note="4 pending · 8 submitted" />
      <Metric icon="₹" label="Pending Payments" value={5} note={`${money(284500)} outstanding`} />
      <Metric icon="◈" label="Screening Alerts" value={metrics.screening} note="1 review · background screening" />
    </div>
    <div className="dashboardGrid">
      <section className="panel chartPanel"><div className="panelTitle"><h3>Transaction Summary</h3><button onClick={()=>onOpen(tx[0])}>View all transactions →</button></div><div className="summaryBody"><div className="donut"><strong>{metrics.transactions}</strong><span>Total</span></div><div className="legend">{[['Export of Goods',9],['Import of Goods',7],['HSS',4],['Import LC',3],['Bank Guarantee',2],['Other',17]].map(([x,n])=><div key={String(x)}><i></i><span>{x}</span><b>{n}</b></div>)}</div></div></section>
      <section className="panel chartPanel"><div className="panelTitle"><h3>Transaction Status</h3></div><div className="bars">{[['Draft',6],['In Progress',14],['Pending Action',8],['With Bank',7],['Completed',7]].map(([x,n])=><div key={String(x)}><div className="bar" style={{height:`${Number(n)*8+25}px`}}></div><b>{n}</b><small>{x}</small></div>)}</div></section>
      <section className="panel quick"><h3>Quick Actions</h3><button onClick={onNew}>＋ <span>Create New Transaction</span> ›</button><button onClick={onUpload}>⇧ <span>Upload Document</span> ›</button><button onClick={onPayment}>₹ <span>Record Payment</span> ›</button><button onClick={()=>onScreen()}>◈ <span>Run Screening</span> ›</button><button>?</button><button>◎ <span>Ask Help Me</span> ›</button></section>
    </div>
    <div className="dashboardGrid lower">
      <section className="panel widePanel"><div className="panelTitle"><h3>Recent Transactions</h3><button onClick={()=>onOpen(tx[0])}>View all →</button></div><div className="tableWrap"><table><thead><tr><th>Txn No.</th><th>Product</th><th>Buyer / Supplier</th><th>Country</th><th>Amount</th><th>Status</th><th>Next Action</th></tr></thead><tbody>{tx.slice(0,5).map(t=><tr key={t.id} onClick={()=>onOpen(t)}><td className="link">{t.reference_no}</td><td>{t.product_type}</td><td>{t.counterparty||'—'}</td><td>{t.country||'—'}</td><td>{money(t.declared_value,t.currency)}</td><td><Badge status={t.status}/></td><td>{t.status==='With Bank'?'Prepare bank docs':t.outstanding_value?'Review outstanding':'—'}</td></tr>)}</tbody></table></div></section>
      <section className="panel notifications"><div className="panelTitle"><h3>Recent Notifications</h3><button>View all →</button></div><Notice title="Screening Review" text="Metro Exports has a possible low-confidence match." tone="red"/><Notice title="Bank Requirement" text="Additional documents requested for EXP-2026-000125." tone="orange"/><Notice title="Action Required" text="Payment not yet recorded for IMP-2026-000098." tone="orange"/><Notice title="Transaction Updated" text="LC-2026-000065 status changed to With Bank." tone="green"/></section>
    </div>
    <section className="panel screeningPreview"><div className="panelTitle"><div><h3>Screening</h3><p>Background screening of relevant parties. You are notified when review is needed.</p></div><button className="primary" onClick={onScreen}>Open Screening</button></div><div className="screenTable"><div className="screenRow header"><span>Date</span><span>Entity</span><span>Type</span><span>Result</span><span>Details</span><span>Action</span></div>{screening.slice(0,4).map(s=><div className="screenRow" key={s.id}><span>{s.date}</span><span><b>{s.entity}</b></span><span>{s.role}</span><span><Badge status={s.result}/></span><span>{s.detail}</span><span className="link">View</span></div>)}</div><div className="infoBox"><b>Screening works in the background.</b><span>It supports customer recordkeeping and review. A screening result is not by itself a final regulatory or bank decision.</span></div></section>
  </>
}

function Metric({icon,label,value,note}:{icon:string,label:string,value:any,note:string}){return <div className="metric"><div className="metricIcon">{icon}</div><div><span>{label}</span><strong>{value}</strong><small>{note}</small></div></div>}
function Badge({status}:{status:string}){const cls=status.toLowerCase().replace(/\s+/g,'-');return <span className={`badge ${cls}`}>{status}</span>}
function Notice({title,text,tone}:{title:string,text:string,tone:string}){return <div className="noticeItem"><i className={tone}></i><div><b>{title}</b><span>{text}</span></div></div>}

function Transactions({tx,products,filter,setFilter,onOpen}:{tx:Tx[],products:Product[],filter:string,setFilter:(s:string)=>void,onOpen:(t:Tx)=>void}){return <><div className="pageHead"><div><h2>Transactions</h2><p>One place for your trade records, documents, payments and requirements.</p></div></div><section className="panel"><div className="filters"><button className={filter==='All'?'active':''} onClick={()=>setFilter('All')}>All</button>{products.map(p=><button className={filter===p.name?'active':''} onClick={()=>setFilter(p.name)} key={p.code}>{p.name}</button>)}</div><div className="tableWrap"><table><thead><tr><th>Reference</th><th>Product</th><th>Counterparty</th><th>Country</th><th>Value</th><th>Status</th><th>Outstanding</th></tr></thead><tbody>{tx.map(t=><tr key={t.id} onClick={()=>onOpen(t)}><td className="link">{t.reference_no}</td><td>{t.product_type}</td><td>{t.counterparty||'—'}</td><td>{t.country||'—'}</td><td>{money(t.declared_value,t.currency)}</td><td><Badge status={t.status}/></td><td>{money(t.outstanding_value,t.currency)}</td></tr>)}</tbody></table></div></section></>}

function ScreeningPage({rows,filter,setFilter,run}:{rows:Screening[],filter:string,setFilter:(s:string)=>void,run:(e:string)=>void}){return <><div className="pageHead"><div><h2>Screening</h2><p>Review screening activity for customers, buyers, suppliers and counterparties.</p></div><button className="primary" onClick={()=>run('Demo Counterparty')}>＋ Run Screening</button></div><section className="metricGrid screeningMetrics"><Metric icon="✓" label="Clear" value={rows.filter(x=>x.result==='Clear').length} note="No potential match found"/><Metric icon="!" label="Review" value={rows.filter(x=>x.result==='Review').length} note="Human review required"/><Metric icon="⚠" label="Alerts" value={rows.filter(x=>x.result==='Alert').length} note="Needs attention"/><Metric icon="↻" label="Screened" value={rows.length} note="Recent screening records"/></section><section className="panel"><div className="screenTabs">{['All','Clear','Review','Alert'].map(x=><button className={filter===x?'active':''} onClick={()=>setFilter(x)} key={x}>{x}</button>)}</div><div className="tableWrap"><table><thead><tr><th>Date</th><th>Entity</th><th>Role</th><th>Country</th><th>Result</th><th>Details</th><th>Source</th><th>Action</th></tr></thead><tbody>{rows.map(s=><tr key={s.id}><td>{s.date}</td><td><b>{s.entity}</b><small>{s.id}</small></td><td>{s.role}</td><td>{s.country}</td><td><Badge status={s.result}/></td><td>{s.detail}</td><td>{s.source}</td><td><button className="textButton">View</button></td></tr>)}</tbody></table></div></section><div className="twoCols"><section className="panel"><h3>How screening works</h3><p>Screening can run when a relevant party is created or updated, and can be repeated when required. Results are recorded with date, source and review status.</p><div className="infoBox"><b>Important:</b><span>A potential match is a review signal, not a final determination. The customer and appropriate bank/specialist decide the next step.</span></div></section><section className="panel"><h3>Screening settings</h3><div className="settingRow"><span>Screen new counterparties</span><b>On</b></div><div className="settingRow"><span>Rescreen on change</span><b>On</b></div><div className="settingRow"><span>Background screening</span><b>On</b></div><div className="settingRow"><span>Adverse media</span><b>Optional</b></div></section></div></>}

function Transaction360({tx,onBack,onEdit,onRemark}:{tx:Tx,onBack:()=>void,onEdit:()=>void,onRemark:(remark:string)=>void}){
  const [remark,setRemark]=useState(tx.remarks||'')
  return <><div className="pageHead"><div><button className="back" onClick={onBack}>← Transactions</button><h2>{tx.reference_no}</h2><p>{tx.product_type} · {tx.counterparty||'Counterparty not recorded'} · {tx.country||'Country not recorded'}</p></div><div className="statusActions"><Badge status={tx.status}/><button className="primary" onClick={onEdit}>Edit details</button><button onClick={()=>document.getElementById('tx-remarks')?.scrollIntoView({behavior:'smooth'})}>Add remark</button><button>Help Me</button></div></div>
  <section className="summaryStrip"><div><span>Transaction Value</span><b>{money(tx.declared_value,tx.currency)}</b></div><div><span>Invoice Value</span><b>{money(tx.invoice_value,tx.currency)}</b></div><div><span>Realized</span><b>{money(tx.realized_value,tx.currency)}</b></div><div><span>Outstanding</span><b>{money(tx.outstanding_value,tx.currency)}</b></div><div><span>Difference</span><b>{money(tx.difference,tx.currency)}</b></div></section>
  <div className="twoCols"><section className="panel"><h3>What needs attention?</h3>{tx.outstanding_value>0?<Notice title="Outstanding amount" text={`${money(tx.outstanding_value,tx.currency)} is not yet matched to a recorded realization.`} tone="orange"/>:<Notice title="No outstanding payment" text="There is no recorded outstanding amount in the current customer record." tone="green"/>}<h3>Next steps</h3><div className="step"><b>1</b><div><strong>Keep the transaction record updated</strong><span>Edit details whenever the customer record changes.</span></div></div><div className="step"><b>2</b><div><strong>Keep supporting evidence together</strong><span>Add documents, bank references and remarks as they become available.</span></div></div><div className="step"><b>3</b><div><strong>Review and reconcile</strong><span>Match invoice, payment and transaction values and review differences.</span></div></div></section>
  <section className="panel"><h3>Transaction record</h3><div className="detailList"><div><span>Product</span><b>{tx.product_type}</b></div><div><span>Currency</span><b>{tx.currency}</b></div><div><span>Counterparty</span><b>{tx.counterparty||'—'}</b></div><div><span>Country</span><b>{tx.country||'—'}</b></div><div><span>Description</span><b>{tx.description||'—'}</b></div><div><span>Customer remarks</span><b>{tx.remarks||'—'}</b></div></div><div className="infoBox"><b>Regulatory intelligence</b><span>Relevant rules can be surfaced here only when they affect a document, deadline, declaration or exception.</span></div></section></div>
  <section className="panel remarkPanel" id="tx-remarks" style={{marginTop:14}}><div className="panelTitle"><div><h3>Customer remarks</h3><p>Use this for context, bank follow-up notes, internal explanations or items you want to remember. It is not a bank submission.</p></div><span className="safeTag">Private customer record</span></div><textarea value={remark} onChange={e=>setRemark(e.target.value)} placeholder="Add a remark…" maxLength={2000}/><div className="remarkFooter"><small>{remark.length}/2000</small><button className="primary" onClick={()=>onRemark(remark.trim())}>Save remark</button></div></section></>}

function EditTransactionModal({tx,onCancel,onSave}:{tx:Tx,onCancel:()=>void,onSave:(tx:Tx)=>void}){
  const [form,setForm]=useState(tx)
  const set=(k:keyof Tx,v:any)=>setForm(p=>({...p,[k]:v}))
  return <div className="modal"><div className="modalCard"><div className="modalHead"><div><h2>Edit transaction</h2><p>Customer-controlled record. Changes do not submit anything to the bank.</p></div><button onClick={onCancel}>×</button></div><div className="twoCols"><label>Reference number<input value={form.reference_no} onChange={e=>set('reference_no',e.target.value)} maxLength={80}/></label><label>Status<select value={form.status} onChange={e=>set('status',e.target.value)}><option>Draft</option><option>In Progress</option><option>Attention Required</option><option>With Bank</option><option>Completed</option></select></label><label>Currency<input value={form.currency} onChange={e=>set('currency',e.target.value.toUpperCase())} maxLength={3}/></label><label>Transaction value<input type="number" min="0" value={form.declared_value} onChange={e=>set('declared_value',Number(e.target.value))}/></label><label>Counterparty<input value={form.counterparty} onChange={e=>set('counterparty',e.target.value)} maxLength={200}/></label><label>Country<input value={form.country} onChange={e=>set('country',e.target.value)} maxLength={100}/></label><label>Description<textarea value={form.description} onChange={e=>set('description',e.target.value)} maxLength={1000}/></label><label>Customer remarks<textarea value={form.remarks||''} onChange={e=>set('remarks',e.target.value)} maxLength={2000}/></label></div><div className="infoBox"><b>Safe editing</b><span>Financial realization values are not directly overwritten here. Payment and realization records remain the source for reconciliation.</span></div><div className="modalActions"><button onClick={onCancel}>Cancel</button><button className="primary" onClick={()=>onSave({...form,reference_no:form.reference_no.trim(),currency:form.currency.trim().toUpperCase()})}>Save changes</button></div></div></div>}

function SecuritySettings({locked,onLock,auditLog}:{locked:boolean,onLock:()=>void,auditLog:AuditEvent[]}){return <><div className="pageHead"><div><h2>Security & Settings</h2><p>Clear controls for customers, with security safeguards separated from trade decisions.</p></div><button className="primary" onClick={onLock}>{locked?'Locked':'Lock workspace'}</button></div><div className="securityGrid"><section className="panel"><h3>Current security posture</h3><div className="securityItem"><b>Server authentication</b><span>Production backend supports authenticated requests; replace prototype tokens with OIDC/OAuth/JWT before live use.</span></div><div className="securityItem"><b>Transport protection</b><span>Security headers and HSTS are enabled for production mode. Use HTTPS end-to-end.</span></div><div className="securityItem"><b>Rate limiting</b><span>API requests are rate-limited to reduce abuse and brute-force attempts.</span></div><div className="securityItem"><b>Tenant separation</b><span>Authenticated principals carry company scope; production handlers must enforce it on every read/write.</span></div><div className="securityItem"><b>Audit trail</b><span>Customer edits are recorded in this prototype activity history; production audit events must be server-side and tamper-resistant.</span></div><button className="secondaryAction" onClick={onLock}>Lock visible workspace</button></section><section className="panel"><h3>Customer-friendly safeguards</h3><div className="securityItem"><b>Edit instead of re-enter</b><span>Customers can correct transaction details without creating duplicate records.</span></div><div className="securityItem"><b>Remarks everywhere</b><span>Context can be captured without turning remarks into mandatory regulatory fields.</span></div><div className="securityItem"><b>Confirmation before risky changes</b><span>Future production workflow should require confirmation and reason for sensitive financial corrections.</span></div><div className="securityItem"><b>No automatic bank submission</b><span>The application remains a customer-side preparation and recordkeeping tool.</span></div></section></div><section className="panel" style={{marginTop:14}}><div className="panelTitle"><h3>Recent activity</h3><span className="safeTag">Customer audit preview</span></div>{auditLog.length===0?<p className="muted">No edits recorded in this session yet.</p>:<div className="tableWrap"><table><thead><tr><th>Time</th><th>Area</th><th>Action</th><th>Reference</th><th>Detail</th></tr></thead><tbody>{auditLog.slice(0,20).map(a=><tr key={a.id}><td>{a.at}</td><td>{a.area}</td><td>{a.action}</td><td>{a.reference}</td><td>{a.detail}</td></tr>)}</tbody></table></div>}</section></>}

function DeclarationsPage({tx,onNotice}:{tx:Tx[],onNotice:(s:string)=>void}){
  const [txId,setTxId]=useState(tx[0]?.id||'')
  const [invoiceNo,setInvoiceNo]=useState('')
  const [invoiceDate,setInvoiceDate]=useState('')
  const [invoiceMonth,setInvoiceMonth]=useState('')
  const [value,setValue]=useState('')
  const [authority,setAuthority]=useState('')
  const [status,setStatus]=useState('Customer review pending')
  const [submittedDate,setSubmittedDate]=useState('')
  const [submissionRef,setSubmissionRef]=useState('')
  const [records,setRecords]=useState<any[]>([])
  const current=tx.find(t=>t.id===txId)
  const isExport=!!current && (current.product_code.startsWith('EXP') || current.product_type.toLowerCase().includes('export'))
  const declarationType=isExport?'EDF':'IMPORT_DECLARATION'
  const save=async()=>{
    if(!current){onNotice('Choose a transaction first.');return}
    if(!invoiceNo.trim() || !invoiceDate || !value){onNotice('Enter invoice number, invoice date and declaration value before saving.');return}
    const record={transaction_id:current.id,declaration_type:declarationType,version:'2026.1',status,fields:{invoice_number:invoiceNo,invoice_date:invoiceDate,invoice_month:invoiceMonth||invoiceDate.slice(0,7),declaration_value:value,currency:current.currency,specified_authority:authority,submission_date:submittedDate,submission_reference:submissionRef,record_type:isExport?'Export Declaration Form (EDF)':'Import declaration / bank record'}}
    try{
      const res=await fetch(API+'/declarations',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(record)})
      if(!res.ok) throw new Error('API unavailable')
      const saved=await res.json();setRecords(prev=>[saved,...prev]);onNotice('Declaration record saved. This records customer-provided status; it does not submit to the AD bank.')
    }catch{
      setRecords(prev=>[{...record,id:`LOCAL-${Date.now()}`},...prev]);onNotice('Saved in this browser demo only. Start the backend to persist the declaration record.')
    }
  }
  return <><div className="pageHead"><div><h2>Declarations & Filing Records</h2><p>Keep a record of declaration preparation and submission status. Nothing is submitted to a bank from this screen.</p></div></div>
    <div className="infoBox"><b>Important distinction</b><span>EDF is an export declaration for goods, services and software under the applicable framework. Imports use applicable import declarations and monitoring records—not EDF. Confirm applicability and process with your AD bank.</span></div>
    <section className="panel" style={{marginTop:14}}><h3>Record a declaration</h3>
      <div className="detailList"><div><span>Transaction</span><select value={txId} onChange={e=>setTxId(e.target.value)} style={{maxWidth:'70%',padding:8,border:'1px solid #dce4ee',borderRadius:6}}>{tx.map(t=><option key={t.id} value={t.id}>{t.reference_no} · {t.product_type}</option>)}</select></div><div><span>Record type</span><b>{isExport?'Export Declaration Form (EDF)':'Import declaration / bank record'}</b></div></div>
      <div className="twoCols" style={{marginTop:14}}>
        <label>Invoice number<input value={invoiceNo} onChange={e=>setInvoiceNo(e.target.value)} placeholder="e.g. INV-2026-001" /></label>
        <label>Invoice date<input type="date" value={invoiceDate} onChange={e=>setInvoiceDate(e.target.value)} /></label>
        <label>Invoice month (if applicable)<input type="month" value={invoiceMonth} onChange={e=>setInvoiceMonth(e.target.value)} /></label>
        <label>Declaration value<input type="number" min="0" value={value} onChange={e=>setValue(e.target.value)} placeholder="Enter invoice / declaration value" /></label>
        <label>Specified authority / channel<input value={authority} onChange={e=>setAuthority(e.target.value)} placeholder="Enter only when known" /></label>
        <label>Customer-recorded status<select value={status} onChange={e=>setStatus(e.target.value)}><option>Draft</option><option>Customer review pending</option><option>Prepared</option><option>Submitted to AD Bank</option><option>Accepted / recorded (as confirmed)</option><option>Correction required</option></select></label>
        <label>Submission date<input type="date" value={submittedDate} onChange={e=>setSubmittedDate(e.target.value)} /></label>
        <label>Submission / bank reference<input value={submissionRef} onChange={e=>setSubmissionRef(e.target.value)} placeholder="Optional reference" /></label>
      </div>
      <div style={{display:'flex',justifyContent:'flex-end',marginTop:16}}><button className="primary" onClick={save}>Save declaration record</button></div>
    </section>
    <section className="panel" style={{marginTop:14}}><h3>Records created in this session</h3>{records.length===0?<p style={{color:'#718096',fontSize:12}}>No declaration records created in this session yet.</p>:<div className="tableWrap"><table><thead><tr><th>Transaction</th><th>Type</th><th>Invoice</th><th>Value</th><th>Status</th><th>Submission reference</th></tr></thead><tbody>{records.map((r,i)=><tr key={r.id||i}><td>{tx.find(t=>t.id===r.transaction_id)?.reference_no||r.transaction_id}</td><td>{r.declaration_type}</td><td>{r.fields?.invoice_number}</td><td>{r.fields?.declaration_value} {r.fields?.currency}</td><td><Badge status={r.status}/></td><td>{r.fields?.submission_reference||'—'}</td></tr>)}</tbody></table></div>}</section>
  </>
}

function Placeholder({tab,tx,onOpen}:{tab:string,tx:Tx[],onOpen:(t:Tx)=>void}){return <><div className="pageHead"><div><h2>{tab}</h2><p>Customer-side workspace. Keep records, see requirements and know what to do next.</p></div></div><section className="panel placeholder"><div className="placeholderIcon">{tab==='Help Me'?'?':'◎'}</div><h3>{tab} is the next working screen</h3><p>The navigation and workflow are connected. We will implement the detailed interactions in the next pass and validate them against real customer scenarios.</p><div className="twoCols miniCards"><div><b>Primary action</b><span>{tab==='Documents'?'Upload and classify a document':tab==='Payments & Reconciliation'?'Record and match a payment':tab==='Bank Requirements'?'Record what the bank requested':tab==='Buyer / Supplier Requirements'?'Track external requirements':tab==='Action Center'?'Work through open actions':tab==='Reports'?'View management and transaction reports':'Ask a transaction-aware question'}</span></div><div><b>Test with</b><span>{tx[0]?.reference_no || 'Create a new transaction'} and a realistic customer case</span><button onClick={()=>onOpen(tx[0])}>Open Transaction 360 →</button></div></div></section></>}


function PaymentsWorkspace({tx,records,onAdd}:{tx:Tx[],records:PaymentRecord[],onAdd:(r:PaymentRecord)=>void}){
  const [transactionId,setTransactionId]=useState(tx.find(t=>t.product_code==='EXP_SERVICE')?.id||tx[0]?.id||'')
  const [direction,setDirection]=useState<'Inward'|'Outward'>('Inward')
  const [reference,setReference]=useState(''); const [amount,setAmount]=useState(''); const [date,setDate]=useState(''); const [invoiceRef,setInvoiceRef]=useState('')
  type Invoice={id:string;txId:string;number:string;date:string;value:number;currency:string;lodged:boolean;bankRef:string;lodgedDate:string;remarks:string;lodgementRemarks:string}
  type Realization={id:string;invoiceId:string;irm:string;date:string;amount:number;thirdParty:string;note:string}
  const [invoices,setInvoices]=useState<Invoice[]>([]); const [realizations,setRealizations]=useState<Realization[]>([])
  const [invNo,setInvNo]=useState(''); const [invDate,setInvDate]=useState(''); const [invValue,setInvValue]=useState(''); const [invRemarks,setInvRemarks]=useState('');
  const [lodgeInvoice,setLodgeInvoice]=useState(''); const [bankRef,setBankRef]=useState(''); const [lodgeDate,setLodgeDate]=useState(''); const [lodgeRemarks,setLodgeRemarks]=useState('');
  const [realInvoice,setRealInvoice]=useState(''); const [irm,setIrm]=useState(''); const [realDate,setRealDate]=useState(''); const [realAmount,setRealAmount]=useState(''); const [thirdParty,setThirdParty]=useState(''); const [realNote,setRealNote]=useState('');
  const selected=tx.find(t=>t.id===transactionId); const serviceTx=selected?.product_code==='EXP_SERVICE';
  const outstanding=(i:Invoice)=>i.value-realizations.filter(r=>r.invoiceId===i.id).reduce((a,r)=>a+r.amount,0)
  const addInvoice=()=>{if(!selected||!invNo.trim()||!invDate||Number(invValue)<=0)return;const i:Invoice={id:`INV-${Date.now()}`,txId:selected.id,number:invNo.trim(),date:invDate,value:Number(invValue),currency:selected.currency,lodged:false,bankRef:'',lodgedDate:'',remarks:invRemarks.trim(),lodgementRemarks:''};setInvoices(p=>[i,...p]);setInvNo('');setInvDate('');setInvValue('');setInvRemarks('')}
  const lodge=()=>{if(!lodgeInvoice||!bankRef.trim()||!lodgeDate)return;setInvoices(p=>p.map(i=>i.id===lodgeInvoice?{...i,lodged:true,bankRef:bankRef.trim(),lodgedDate:lodgeDate,lodgementRemarks:lodgeRemarks.trim()}:i));setBankRef('');setLodgeDate('');setLodgeRemarks('')}
  const realize=()=>{const i=invoices.find(x=>x.id===realInvoice);const n=Number(realAmount);if(!i||!irm.trim()||!realDate||n<=0||n>outstanding(i))return;setRealizations(p=>[{id:`REAL-${Date.now()}`,invoiceId:i.id,irm:irm.trim(),date:realDate,amount:n,thirdParty:thirdParty.trim(),note:realNote.trim()},...p]);setIrm('');setRealDate('');setRealAmount('');setThirdParty('');setRealNote('')}
  const save=()=>{if(!selected||!reference.trim()||!amount||Number(amount)<=0||!date)return;onAdd({id:`PAY-${Date.now()}`,transactionId,direction,reference,amount:Number(amount),date,currency:selected.currency,invoiceRef});setReference('');setAmount('');setDate('');setInvoiceRef('')}
  return <><div className="pageHead"><div><h2>Payments & Reconciliation</h2><p>Record invoices, separate bank lodgement from realization, and review outstanding balances.</p></div></div><div className="infoBox"><b>Prototype boundary</b><span>Service invoice, lodgement and realization records are held in this browser session only. “Lodged” means customer-recorded submission/reference, not verified bank acceptance or official filing. Confirm applicable declaration and bank process with your AD bank.</span></div>
  <section className="panel" style={{marginTop:14}}><h3>Service / software export lifecycle</h3><label style={{display:'block',margin:'12px 0'}}>Transaction<select value={transactionId} onChange={e=>setTransactionId(e.target.value)}>{tx.map(t=><option key={t.id} value={t.id}>{t.reference_no} · {t.product_type}</option>)}</select></label>
  {serviceTx?<><div className="twoCols"><div><h4>1. Register service invoice</h4><label>Invoice number<input value={invNo} onChange={e=>setInvNo(e.target.value)} placeholder="e.g. SVC-001"/></label><label>Invoice date<input type="date" value={invDate} onChange={e=>setInvDate(e.target.value)}/></label><label>Invoice value ({selected?.currency})<input type="number" min="0.01" value={invValue} onChange={e=>setInvValue(e.target.value)}/></label><label>Customer remarks (optional)<textarea value={invRemarks} onChange={e=>setInvRemarks(e.target.value)} placeholder="Internal note or context"/></label><button className="primary" style={{marginTop:10}} disabled={!invNo.trim()||!invDate||Number(invValue)<=0} onClick={addInvoice}>Add invoice</button></div><div><h4>2. Record lodge with AD bank</h4><label>Invoice<select value={lodgeInvoice} onChange={e=>setLodgeInvoice(e.target.value)}><option value="">Select invoice</option>{invoices.filter(i=>i.txId===transactionId&&!i.lodged).map(i=><option key={i.id} value={i.id}>{i.number}</option>)}</select></label><label>Bank acknowledgement / reference<input value={bankRef} onChange={e=>setBankRef(e.target.value)} placeholder="Customer-entered bank reference"/></label><label>Lodgement date<input type="date" value={lodgeDate} onChange={e=>setLodgeDate(e.target.value)}/></label><label>Lodgement remarks (optional)<textarea value={lodgeRemarks} onChange={e=>setLodgeRemarks(e.target.value)} placeholder="Bank response, pending point, follow-up"/></label><button className="primary" style={{marginTop:10}} disabled={!lodgeInvoice||!bankRef.trim()||!lodgeDate} onClick={lodge}>Record lodgement</button></div></div>
  <div style={{marginTop:18}}><h4>3. Record realization / IRM and allocate to invoice</h4><div className="twoCols"><label>Invoice<select value={realInvoice} onChange={e=>setRealInvoice(e.target.value)}><option value="">Select invoice</option>{invoices.filter(i=>i.txId===transactionId&&outstanding(i)>0).map(i=><option key={i.id} value={i.id}>{i.number} · outstanding {money(outstanding(i),i.currency)}</option>)}</select></label><label>IRM / bank remittance reference<input value={irm} onChange={e=>setIrm(e.target.value)} placeholder="Bank advice / IRM reference"/></label><label>Realization date<input type="date" value={realDate} onChange={e=>setRealDate(e.target.value)}/></label><label>Allocated amount<input type="number" min="0.01" value={realAmount} onChange={e=>setRealAmount(e.target.value)}/></label><label>Third-party payer (optional)<input value={thirdParty} onChange={e=>setThirdParty(e.target.value)} placeholder="If payer differs from buyer"/></label><label>Difference / notes (optional)<input value={realNote} onChange={e=>setRealNote(e.target.value)} placeholder="Short payment, charges, explanation…"/></label></div><button className="primary" style={{marginTop:10}} disabled={!realInvoice||!irm.trim()||!realDate||Number(realAmount)<=0||Number(realAmount)>(invoices.find(i=>i.id===realInvoice)?outstanding(invoices.find(i=>i.id===realInvoice)!):0)} onClick={realize}>Allocate realization</button></div>
  <div className="tableWrap" style={{marginTop:18}}><table><thead><tr><th>Invoice</th><th>Invoice value</th><th>Lodgement</th><th>Realized</th><th>Outstanding</th><th>Reconciliation</th><th>Remarks / action</th></tr></thead><tbody>{invoices.filter(i=>i.txId===transactionId).map(i=>{const got=i.value-outstanding(i);return <tr key={i.id}><td>{i.number}<small style={{display:'block'}}>{i.date}</small></td><td>{money(i.value,i.currency)}</td><td>{i.lodged?`Recorded · ${i.bankRef}`:'Not recorded'}</td><td>{money(got,i.currency)}</td><td>{money(outstanding(i),i.currency)}</td><td>{outstanding(i)===0?'Fully allocated':got>0?'Partly realized':'Awaiting realization'}</td><td><small>{[i.remarks,i.lodgementRemarks].filter(Boolean).join(' · ')||'—'}</small><button type="button" onClick={()=>{const number=window.prompt('Invoice number',i.number);if(number===null)return;const date=window.prompt('Invoice date (YYYY-MM-DD)',i.date);if(date===null)return;const value=window.prompt('Invoice value',String(i.value));if(value===null)return;const remarks=window.prompt('Customer remarks',i.remarks);if(remarks===null)return;const n=Number(value);if(!number.trim()||!date||!Number.isFinite(n)||n<=0||n< i.value-outstanding(i)){window.alert('Invalid update. Value cannot be below the amount already realized.');return;}setInvoices(prev=>prev.map(x=>x.id===i.id?{...x,number:number.trim(),date,value:n,remarks}:x));}}>Edit invoice</button></td></tr>})}{invoices.filter(i=>i.txId===transactionId).length===0&&<tr><td colSpan={7}>No service invoices recorded for this transaction yet.</td></tr>}</tbody></table></div>
  {realizations.some(r=>invoices.some(i=>i.id===r.invoiceId&&i.txId===transactionId))&&<div className="tableWrap" style={{marginTop:12}}><h4>Realization allocation history</h4><table><thead><tr><th>Date</th><th>Invoice</th><th>IRM reference</th><th>Amount</th><th>Payer / notes</th></tr></thead><tbody>{realizations.filter(r=>invoices.some(i=>i.id===r.invoiceId&&i.txId===transactionId)).map(r=><tr key={r.id}><td>{r.date}</td><td>{invoices.find(i=>i.id===r.invoiceId)?.number}</td><td>{r.irm}</td><td>{money(r.amount,invoices.find(i=>i.id===r.invoiceId)?.currency)}</td><td>{[r.thirdParty,r.note].filter(Boolean).join(' · ')||'—'}</td></tr>)}</tbody></table></div>}</>:<p style={{color:'#718096'}}>Select a Software / Service Export transaction to use the invoice → lodge → realization workflow. The general payment recorder remains available below.</p>}</section>
  <section className="panel" style={{marginTop:14}}><h3>General payment record</h3><div className="twoCols" style={{marginTop:12}}><label>Direction<select value={direction} onChange={e=>setDirection(e.target.value as 'Inward'|'Outward')}><option>Inward</option><option>Outward</option></select></label><label>Payment / bank reference<input value={reference} onChange={e=>setReference(e.target.value)} placeholder="e.g. bank reference"/></label><label>Amount<input type="number" min="0.01" value={amount} onChange={e=>setAmount(e.target.value)} placeholder="Enter amount"/></label><label>Payment date<input type="date" value={date} onChange={e=>setDate(e.target.value)}/></label><label>Invoice reference (optional)<input value={invoiceRef} onChange={e=>setInvoiceRef(e.target.value)} placeholder="Invoice to associate"/></label></div><div style={{display:'flex',justifyContent:'flex-end',marginTop:14}}><button className="primary" onClick={save} disabled={!selected||!reference.trim()||!amount||Number(amount)<=0||!date}>Record payment</button></div></section>
  <section className="panel" style={{marginTop:14}}><h3>Recorded general payments</h3>{records.length===0?<p style={{color:'#718096',fontSize:12}}>No general payments recorded in this session.</p>:<div className="tableWrap"><table><thead><tr><th>Date</th><th>Transaction</th><th>Direction</th><th>Reference</th><th>Invoice</th><th>Amount</th></tr></thead><tbody>{records.map(r=><tr key={r.id}><td>{r.date}</td><td>{tx.find(t=>t.id===r.transactionId)?.reference_no||r.transactionId}</td><td>{r.direction}</td><td>{r.reference}</td><td>{r.invoiceRef||'—'}</td><td>{money(r.amount,r.currency)}</td></tr>)}</tbody></table></div>}</section></>}
