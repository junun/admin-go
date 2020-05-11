import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Row, Col, 
  Button, Popconfirm, Icon, message, Radio, Select, Switch,
  Transfer, Steps, Tag} from "antd";
import {connect} from "dva";
import Link from 'umi/link';
import ShowWebSocketMessageTemplate from '@/components/Report/ShowWebSocketMessageTemplate';
import {timeTrans, compareArray, hasPermission} from "@/utils/globalTools";
import Setup1 from './Ext1Setup1';
import Setup2 from './Ext1Setup2';
import Setup3 from './Ext1Setup3'; 
import styles from './index.module.css';

const FormItem = Form.Item;
const Option = Select.Option;
const RadioButton = Radio.Button;
const RadioGroup = Radio.Group;

@connect(({ loading, config, host }) => {
  return {
    projectsList: config.projectsList,
    projectsLoading: loading.effects['config/getProject'],
    projectsLen: config.projectsLen,
    configEnvList: config.configEnvList,
    appTypeList: config.appTypeList,
    hostListByAppId: host.hostListByAppId,
    deployTempleList:config.deployTempleList,
  }
})

class ProjectListPage extends React.Component {
  constructor(props) {
    super(props)
    this.handler = this.handler.bind(this)
  }

  state = {
    visible: false,
    showTemplate: false,
    editCacheData: {},
    body: '',
    id: 0,
    stepsVisible: false,
    pageNum: 0,
    info:{},
    loadDeployTempleList: 0,
    secondTable: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'config/getProject',
      payload: {
        page: 1,
        pageSize: 10, 
      }
    });
    dispatch({
      type: 'config/getConfigEnv',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    });
    dispatch({
      type: 'config/getAppType',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    });
  };

  handler = (values) => {
    this.setState({
      info: values,
      pageNum: this.state.pageNum + 1,
    })
  }

  getDeployExtend = (values) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'config/getDeployExtend',
      payload: {
        id: values
      }
    }).then(()=> {
      var tmp    = this.state.secondTable
      tmp[values]= this.props.deployTempleList
      this.setState({
        secondTable: tmp,
        loadDeployTempleList: values,
      })
    });
  }

  prehandler = (values) => {
    this.setState({
      pageNum: this.state.pageNum - 1,
    })
  }

  showEnvAddModal = () => {
    this.setState({
      isOk : false,
      editCacheData: {},
      visible: true,
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
      stepsVisible: false,
      info: {},
      deployTempleList: 0,
      pageNum: 0,
    });
  };

  // 关联角色是为了后期 deploy 权限检查
  // 数组比较
  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        if (Object.keys(obj).length) {
          values.Active      =  values.Active ? 1 : 0
          values.EnableSync  =  values.EnableSync ? 1 : 0
          values.DeployType  =  values.DeployType ? 1 : 0
          if (
            obj.Name          === values.Name &&
            obj.DeployType    === values.DeployType &&
            obj.EnableSync    === values.EnableSync &&
            obj.Active        === values.Active &&
            obj.Desc          === values.Desc &&
            obj.Tid           === values.Tid &&
            obj.EnvId         === values.EnvId &&
            obj.RepoUrl       === values.RepoUrl
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'config/configProjectEdit',
              payload: values,
            }).then(()=>{
              this.setState({ visible: false });
            });
          }
        } else {
          values.Active      = values.Active ? 1 : 0
          values.EnableSync  = values.EnableSync ? 1 : 0
          values.DeployType  =  values.DeployType ? 1 : 0
          dispatch({
            type: 'config/configProjectAdd',
            payload: values,
          }).then(()=>{
            this.setState({ visible: false })
          });
        }
      }
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'config/configProjectDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  // Popconfirm 取消事件
  cancel = () => {
  };


  //显示编辑界面
  handleEdit = (values, page) => {
    values.title =  '编辑项目-' + values.Name;
    this.setState({ 
      visible: true ,
      editCacheData: values
    });
  };

  // 显示项目发布配置模板，用于发布申请
  handleAppDeployConfig = (values) => {
    this.setState({ 
      ext1FromVisible: !this.state.ext1FromVisible,
      id: values.id,
      editCacheData: values,
    });
  };

  showStepsAddModal = (id, info, isClone) => {
    if (info === undefined) {
      var tmp = this.state.info
      tmp["Aid"] = id
    } else {
      var tmp = info
      tmp["Aid"] = id
    }

    isClone && delete tmp.Dtid

    this.setState({
      stepsVisible: true,
      info: tmp,
    });
  };

  handleSync = (values) => {
    this.setState({ 
      showTemplate: true ,
      id: values,
    });
  };

  onChange = (e) => {
  };

  onCheckChange = (checked) => {
  };

  switchTemplate = () => {
    this.setState({ 
      showTemplate: !this.state.showTemplate,
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'task/getProject',
      payload: {
        page: page,
        pageSize: 10, 
      }
    });
  };

  nextPage = () => {
    this.setState({ 
      
    });
  };

  handleDeployDelete = (id, dtid) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'config/deployExtendDel',
      payload: {
        id: id,
        tid: dtid,
      }
    });
  };

  columns = [
    {
      title: '项目名',
      dataIndex: 'Name',
    },{
      title: '发布环境',
      dataIndex: 'EnvId',
      'render': EnvId => this.props.configEnvList.map(x => {
        if (EnvId == x.id) {
          return x.Name
        }
      })
    },{
      title: '项目类型',
      dataIndex: 'Tid',
      'render': Tid => this.props.appTypeList.map(x => {
        if (Tid == x.id) {
          return x.Name
        }
      })
    },{
      title: '发布类型',
      dataIndex: 'DeployType',
      'render': DeployType => DeployType && '自定义发布' || '目录copy',
    },{
      title: '状态',
      dataIndex: 'Active',
      'render': Active => Active && '启用' || '禁用',
    }, {
      title: '备注',
      dataIndex: 'Desc',
      ellipsis: true,
    }, {
      title: '操作',
      key: 'action',
      render: (text, record) => (
        <span>
          {
            hasPermission('config-app-edit') && 
            <a onClick={()=>{this.handleEdit(record)}}>
              <Icon type="edit"/>编辑
            </a>
          }
          <Divider type="vertical" />
          {  
            hasPermission('config-app-set') && 
            <a onClick={()=>{this.showStepsAddModal(record.id)}}>
              <Icon type="save"/>新建发布模板
            </a>
          }
          <Divider type="vertical" />
          {  
            record.EnableSync == 1
            &&
            hasPermission('config-app-set') && 
            <Link to={`/config/value?aid=${record.id}`}>
              <Icon type="setting"/>初始化变量
            </Link> 
          }
          {
            record.EnableSync == 1 &&
            <Divider type="vertical" />
          }
          {  
            record.EnableSync == 1 &&
            hasPermission('config-app-init') && 
            <Popconfirm title="你确定要操作吗?"
              onConfirm={()=>{this.handleSync(record.id)}} onCancel={()=>{this.cancel()}}>
              <a title="初始化" ><Icon type="redo" />初始化</a>
            </Popconfirm>
          }
          {
            record.EnableSync == 1 &&
            <Divider type="vertical" />
          }
          {
            hasPermission('config-app-del') && 
            <Popconfirm title="你确定要删除吗?"  
              onConfirm={()=>{this.deleteRecord(record.id)}} 
              onCancel={()=>{this.cancel()}}
            >
              <a title="删除" >
                 <Icon type="delete" />删除
              </a>
            </Popconfirm>
          }
        </span>
      ),
    },
  ];
  
  render() {
    const {visible, showTemplate, editCacheData, stepsVisible, pageNum,
      body, info, loadDeployTempleList, secondTable } = this.state;

    const {configEnvList, appTypeList, projectsList, projectsLoading, 
      projectsLen, hostsList, projectTargetKeys, imageList, hostListByAppId,
      deployTempleList,
      form: { getFieldDecorator} } = this.props;

    const addvar = <Button type="" onClick={this.showEnvAddModal} >添加</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('config-app-add') && <Col span={10}>{addvar}</Col>}
      </Row>;

    const expandedRowRender = (record) => {
      if (!secondTable[record.id]) {
        this.getDeployExtend(record.id)
      }

      // if (loadDeployTempleList!=record.id) {
      //   this.getDeployExtend(record.id)
      // }

      const columns = [
      {
        title: '模板名字',
        dataIndex: 'TemplateName',
      }, {
        title: '发布审核',
        dataIndex: 'EnableCheck',
        render: value => value ? <Tag color="green">开启</Tag> : <Tag color="red">关闭</Tag>
      }, {
        title: '操作',
        render: info => (
          <span>
            { 
              hasPermission('config-app-edit') &&  
              <a onClick={()=>{this.showStepsAddModal(record.id, info)}}>编辑</a>
            }
            { 
              hasPermission('config-app-edit') && 
              <Divider type="vertical"/>
            }
            { 
              hasPermission('config-app-edit') && 
               <a onClick={()=>{this.showStepsAddModal(record.id, info, true)}}>克隆配置</a>
            }
            { 
              hasPermission('config-app-edit') && 
              <Divider type="vertical"/>
            }
            {
              hasPermission('config-app-del') && 
              <Popconfirm title="你确定要删除吗?"  
                onConfirm={()=>{this.handleDeployDelete(record.id, info.Dtid)}} 
                onCancel={()=>{this.cancel()}}
              >
                <a title="删除" >
                  删除
                </a>
              </Popconfirm>
            }
          </span>
        )
      }];

      return <Table
        rowKey="Dtid"
        loading={!secondTable[record.id]}
        columns={columns}
        dataSource={secondTable[record.id]}
        pagination={false} />
    };

    return (
      <div>
        <Modal
          title= { editCacheData.title || "新建项目" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
          style={{ top: 20 }} >
          <Form>
            <FormItem label="项目名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="项目类型">
              {getFieldDecorator('Tid', {
                initialValue: editCacheData.Tid || 'Please select' ,
                rules: [{ required: true }],
              })(
                <Select
                  placeholder="Please select"
                  style={{ width: '100%' }}
                >
                  {appTypeList.map(x => <Select.Option key={x.id} value={x.id}>{x.Name}</Select.Option>)}
                </Select>
              )}
            </FormItem>
            <FormItem label="项目所属环境">
              <Col span={16}>
                {getFieldDecorator('EnvId', {
                  initialValue: editCacheData.EnvId || 'Please select' ,
                  rules: [{ required: true }],
                })(
                  <Select
                    placeholder="Please select"
                    style={{ width: '100%' }}
                  >
                    {configEnvList.map(x => <Select.Option key={x.id} value={x.id}>{x.Name}</Select.Option>)}
                  </Select>
                )}
              </Col>
              <Col span={6} offset={2}>
                <Link to="/config/environment">新建环境</Link>
              </Col>
            </FormItem>
            <FormItem label="是否需要初始化">
              {getFieldDecorator('EnableSync', {
                initialValue: editCacheData.EnableSync && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.EnableSync || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="是否生效">
              {getFieldDecorator('Active', {
                initialValue: editCacheData.Active && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.Active || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="是否需要自定义发布">
              {getFieldDecorator('DeployType', {
                initialValue: editCacheData.DeployType && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.DeployType || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="备注信息">
              {getFieldDecorator('Desc', {
                initialValue: editCacheData.Desc || '',
                rules: [{ required: true }],
              })(
                <Input.TextArea />
              )}
            </FormItem>
          </Form> 
        </Modal>

        { showTemplate && 
          <ShowWebSocketMessageTemplate
            id={this.state.id} 
            onCancel={this.switchTemplate} 
            onOk={body => this.setState({body})}
          />
        }

        <Modal
          visible={stepsVisible}
          width={800}
          maskClosable={false}
          destroyOnClose= "true"
          title= { editCacheData.title || "新建常规发布" }
          onCancel={this.handleCancel}
          footer={null}>
          <Steps current={pageNum} className={styles.steps}>
            <Steps.Step key={0} title="基本配置"/>
            <Steps.Step key={1} title="发布主机"/>
            <Steps.Step key={2} title="任务配置"/>
          </Steps>
          {
            pageNum === 0 && 
            <Setup1 
              configEnvList={configEnvList}
              nextPage={this.handler}
              info={this.state.info}
            />
          }
          {
            pageNum === 1 && 
            <Setup2
              info={info}
              nextPage={this.handler}
              prePage={this.prehandler}
              dispatch={this.props.dispatch}
              hostListByAppId={hostListByAppId}
            />
          }
          {
            pageNum === 2 && 
            <Setup3
              info={info}
              prePage={this.prehandler}
              dispatch={this.props.dispatch}
              handleCancel={this.handleCancel}
            />
          }
        </Modal>

        <Card title="" extra={extra}>
          <Table 
            pagination={{
              showQuickJumper: true,
              total: projectsLen,
              showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
              onChange: this.pageChange
            }}
            columns={this.columns}
            dataSource={projectsList}
            expandedRowRender={expandedRowRender}
            loading={projectsLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(ProjectListPage);


